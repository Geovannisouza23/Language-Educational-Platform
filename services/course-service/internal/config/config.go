package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	LogLevel    string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/course-service")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COURSE_SERVICE")

	// Defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("port", "8080")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("database_url", "postgres://postgres:postgres@postgres:5432/courses_db?sslmode=disable")
	viper.SetDefault("redis_url", "redis:6379")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{
		Environment: viper.GetString("environment"),
		Port:        viper.GetString("port"),
		DatabaseURL: viper.GetString("database_url"),
		RedisURL:    viper.GetString("redis_url"),
		JWTSecret:   viper.GetString("jwt_secret"),
		LogLevel:    viper.GetString("log_level"),
	}

	return config, nil
}
