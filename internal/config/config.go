package config

import (
	"errors"
	. "order_service/internal/logger"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	POSTGRES_HOST           string   `env:"POSTGRES_HOST"`
	POSTGRES_PORT           string   `env:"POSTGRES_PORT"`
	POSTGRES_DB             string   `env:"POSTGRES_DB"`
	POSTGRES_USER           string   `env:"POSTGRES_USER"`
	POSTGRES_PASSWORD       string   `env:"POSTGRES_PASSWORD"`
	APP_IP                  string   `env:"APP_IP"`
	APP_PORT                string   `env:"APP_PORT"`
	KAFKA_BOOTSTRAP_SERVER1 string   `env:"KAFKA_BOOTSTRAP_SERVER1"`
	KAFKA_BOOTSTRAP_SERVER2 string   `env:"KAFKA_BOOTSTRAP_SERVER2"`
	KAFKA_BOOTSTRAP_SERVER3 string   `env:"KAFKA_BOOTSTRAP_SERVER3"`
	KAFKA_TOPIC             string   `env:"KAFKA_TOPIC"`
	KAFKA_CONSUMER_GROUP    string   `env:"KAFKA_CONSUMER_GROUP"`
	KAFKA_BOOTSTRAP_SERVERS []string
}

func Load() (*Config, error) {
	var cfg Config

	// Загрузка переменных среды
	if err := godotenv.Load(); err != nil {
		Logger.Warn("No .env file found")
	}

	// Парсинг конфигурации
	if err := env.Parse(&cfg); err != nil {
		errMsg := "Failed to parse config =>" + err.Error()
		Logger.Warn(errMsg)
		return nil, errors.New(errMsg)
	}

	if !cfg.IsValid() {
		errMsg := "bad configuration data"
		Logger.Warn(errMsg, zap.Any("configData", cfg))
		return nil, errors.New(errMsg)
	}

	cfg.KAFKA_BOOTSTRAP_SERVERS = append(cfg.KAFKA_BOOTSTRAP_SERVERS, cfg.KAFKA_BOOTSTRAP_SERVER1)
	cfg.KAFKA_BOOTSTRAP_SERVERS = append(cfg.KAFKA_BOOTSTRAP_SERVERS, cfg.KAFKA_BOOTSTRAP_SERVER2)
	cfg.KAFKA_BOOTSTRAP_SERVERS = append(cfg.KAFKA_BOOTSTRAP_SERVERS, cfg.KAFKA_BOOTSTRAP_SERVER3)

	Logger.Info("Config successfully loaded", zap.Any("configuration", cfg))
	return &cfg, nil
}

func (cfg *Config) IsValid() bool {
	var res bool = true
	if cfg.POSTGRES_USER == "" || cfg.POSTGRES_PASSWORD == "" ||
		cfg.POSTGRES_DB == "" || cfg.POSTGRES_HOST == "" ||
		cfg.POSTGRES_PORT == "" || cfg.APP_IP == "" ||
		cfg.APP_PORT == "" || cfg.KAFKA_BOOTSTRAP_SERVERS == nil ||
		cfg.KAFKA_CONSUMER_GROUP == "" || cfg.KAFKA_TOPIC == "" {
		res = false
	}
	return res
}
