package main

import (
	"order_service/internal/db"
	"order_service/internal/config"
	"order_service/internal/logger"
	"order_service/internal/rest"
	"go.uber.org/zap"
	"syscall"
	"os/signal"
	"os"
	"context"
	"time"
	"net/http"
	"fmt"
)

func main() {
	// Загружаем логгер
	logger := logger.LoadLogger()

	// Загружаем БД
	db, err := db.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer db.Close()
	
	// Создаем REST сервис
	rest.CreateRestService(db)

	// Запуск сервера
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
	srvAddr := fmt.Sprintf("%v:%v", cfg.APP_IP, cfg.APP_PORT)
	srv := &http.Server{
		Addr: srvAddr,
	}

	go func() {
		logger.Info("Server is running", zap.String("address", srvAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Listen error", zap.Error(err))
		}
	}()

	// Выключение сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server is forced to shut down", zap.Error(err))
	}
	logger.Info("Server closed")
}