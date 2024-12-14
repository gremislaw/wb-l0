package logger

import (
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func LoadLogger() {
	// Создание логгера
	l, err := zap.NewProduction()
	if err != nil {
		zap.Error(err)
	}

	// Инициализация логгера
	initLogger(l)
}

func initLogger(l *zap.Logger) {
	Logger = l
}
