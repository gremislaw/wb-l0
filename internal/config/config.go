package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"errors"
)

var (
	logger *zap.Logger
)

type config struct {
	POSTGRES_HOST     string `env:"POSTGRES_HOST"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT"`
	POSTGRES_DB       string `env:"POSTGRES_DB"`
	POSTGRES_USER     string `env:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD"`
	APP_IP            string `env:"APP_IP"`
	APP_PORT          string `env:"APP_PORT"`
}

func Load() (*config, error) {
	var cfg config

	// Загрузка переменных среды
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	// Парсинг конфигурации
	if err := env.Parse(&cfg); err != nil {
		errMsg := "Failed to parse config =>" + err.Error()
		logger.Warn(errMsg)
		return nil, errors.New(errMsg)
	}
	
	logger.Info("Config successfully loaded", zap.String("database", cfg.POSTGRES_DB), zap.String("app_port", cfg.APP_PORT))
	return &cfg, nil
}

func InitLogger(l *zap.Logger) {
	logger = l
}