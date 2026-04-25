package middleware

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP请求处理时间分布",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP请求总数",
		},
		[]string{"method", "path", "status"},
	)

	requestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP请求大小分布",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP响应大小分布",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	dbOpenConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_open_connections",
			Help: "当前数据库打开的连接数",
		},
	)

	dbInUseConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_in_use_connections",
			Help: "当前正在使用的数据库连接数",
		},
	)

	dbIdleConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_idle_connections",
			Help: "当前空闲的数据库连接数",
		},
	)

	dbWaitCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_wait_count_total",
			Help: "等待数据库连接的总次数",
		},
	)

	dbWaitDuration = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_wait_duration_seconds_total",
			Help: "等待数据库连接的总时长（秒）",
		},
	)

	dbMaxOpenConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_max_open_connections",
			Help: "数据库最大允许打开的连接数",
		},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestTotal)
	prometheus.MustRegister(requestSize)
	prometheus.MustRegister(responseSize)
	prometheus.MustRegister(dbOpenConnections)
	prometheus.MustRegister(dbInUseConnections)
	prometheus.MustRegister(dbIdleConnections)
	prometheus.MustRegister(dbWaitCount)
	prometheus.MustRegister(dbWaitDuration)
	prometheus.MustRegister(dbMaxOpenConnections)
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		requestDuration.WithLabelValues(method, path, status).Observe(duration)
		requestTotal.WithLabelValues(method, path, status).Inc()
	}
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// StartDBMetricsCollector 启动数据库连接池指标采集器
func StartDBMetricsCollector(db *sql.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			stats := db.Stats()
			dbOpenConnections.Set(float64(stats.OpenConnections))
			dbInUseConnections.Set(float64(stats.InUse))
			dbIdleConnections.Set(float64(stats.Idle))
			dbWaitCount.Add(float64(stats.WaitCount))
			dbWaitDuration.Add(stats.WaitDuration.Seconds())
			dbMaxOpenConnections.Set(float64(stats.MaxOpenConnections))
		}
	}()
}
