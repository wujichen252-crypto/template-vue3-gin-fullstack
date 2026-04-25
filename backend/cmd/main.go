package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/internal/handler"
	"template-vue3-gin-fullstack/backend/internal/middleware"
	"template-vue3-gin-fullstack/backend/internal/repository"
	"template-vue3-gin-fullstack/backend/internal/service"
	"template-vue3-gin-fullstack/backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	cfg := config.GetConfig()

	logCore := logger.InitLogger(cfg.Server.Mode)
	defer logCore.Sync()

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	rdb := initRedis(cfg)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, rdb)
	userHandler := handler.NewUserHandler(userSvc, rdb, cfg)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.Metrics())
	r.Use(middleware.Logger(logCore))
	r.Use(middleware.Recovery(logCore, false))
	r.Use(middleware.CORS(cfg.Server.AllowOrigins))
	r.Use(middleware.RateLimiter(rdb, cfg))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
	})

	r.GET("/metrics", middleware.MetricsHandler())

	r.Static("/swagger", "./swagger")
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.GET("/userinfo", middleware.JWT(cfg.JWT.Secret, rdb), userHandler.GetUserInfo)
			auth.POST("/refresh", userHandler.RefreshToken)
			auth.POST("/logout", middleware.JWT(cfg.JWT.Secret, rdb), userHandler.Logout)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.Timeout.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Timeout.IdleTimeout) * time.Second,
	}

	go func() {
		log.Printf("服务器启动于 http://localhost:%s", cfg.Server.Port)
		log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
	if rdb != nil {
		rdb.Close()
	}

	log.Println("服务器已退出")
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	maxOpenConns := cfg.Database.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}
	maxIdleConns := cfg.Database.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	connMaxLifetime := cfg.Database.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 3600
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return db, nil
}

func initRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize:     100,
		MinIdleConns: 10,
	})

	return rdb
}