package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string
	Mode         string
	AllowOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret        string
	AccessExpire  int
	RefreshExpire int
}

var globalConfig *Config

func InitConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("配置文件读取失败: %w", err)
		}
	}

	if err := validateConfig(); err != nil {
		return err
	}

	globalConfig = &Config{
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			Mode:         viper.GetString("SERVER_MODE"),
			AllowOrigins: viper.GetStringSlice("ALLOW_ORIGINS"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			AccessExpire:  viper.GetInt("JWT_ACCESS_EXPIRE"),
			RefreshExpire: viper.GetInt("JWT_REFRESH_EXPIRE"),
		},
	}

	return nil
}

func validateConfig() error {
	requiredVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"JWT_SECRET",
	}

	for _, v := range requiredVars {
		if val := viper.GetString(v); val == "" {
			return fmt.Errorf("缺少必需的环境变量: %s", v)
		}
	}

	return nil
}

func GetConfig() *Config {
	if globalConfig == nil {
		if err := InitConfig(); err != nil {
			panic(err)
		}
	}
	return globalConfig
}

func GetJWTConfig() JWTConfig {
	return GetConfig().JWT
}

func GetAccessTokenDuration() time.Duration {
	return time.Duration(GetConfig().JWT.AccessExpire) * time.Hour
}

func GetRefreshTokenDuration() time.Duration {
	return time.Duration(GetConfig().JWT.RefreshExpire) * time.Hour
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}