package logger

import (
	"go.uber.org/zap"
	"order_service/internal/db"
	"order_service/internal/config"
)

func LoadLogger() *zap.Logger {
	// Создание логгера
	logger, err := zap.NewProduction()
	if err != nil {
		zap.Error(err)
	}

	// Инициализация логгера
	db.InitLogger(logger)
	config.InitLogger(logger)

	return logger
}