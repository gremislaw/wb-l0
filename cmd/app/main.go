package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"order_service/internal/config"
	"order_service/internal/db"
	. "order_service/internal/logger"
	"order_service/internal/rest"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Загружаем логгер
	LoadLogger()

	// Загружаем БД
	retries := 5
	db, err := db.Load(retries)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	defer db.Close()

	// Создаем REST сервис
	rest.CreateRestService(db)

	// Запуск сервера
	cfg, err := config.Load()
	if err != nil {
		Logger.Fatal(err.Error())
	}
	srvAddr := fmt.Sprintf("%v:%v", cfg.APP_IP, cfg.APP_PORT)
	srv := &http.Server{
		Addr: srvAddr,
	}

	go func() {
		Logger.Info("Server is running", zap.String("address", srvAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Fatal("Listen error", zap.Error(err))
		}
	}()

	// Выключение сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	Logger.Info("Server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		Logger.Fatal("Server is forced to shut down", zap.Error(err))
	}
	Logger.Info("Server closed")
}
