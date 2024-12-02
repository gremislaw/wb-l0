package db

import (
	"fmt"
	"database/sql"
	"order_service/internal/config"
	"go.uber.org/zap"
	"errors"
)

var (
	logger *zap.Logger
)

func Load() (*sql.DB, error) {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		errMsg := "Failed to load config => " + err.Error()
		logger.Warn(errMsg)
		return nil, errors.New(errMsg)
	}

	// Преобразование конфигурационных данных в DSN
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB, cfg.POSTGRES_HOST, cfg.POSTGRES_PORT)

	// Подключение к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		errMsg := "DB connection error"
		logger.Warn(errMsg, zap.String("DSN", connStr))
		return nil, errors.New(errMsg)
	}

	logger.Info("Database successfully loaded", zap.String("DSN", connStr))
	return db, nil
}

func InitLogger(l *zap.Logger) {
	logger = l
}