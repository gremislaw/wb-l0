package main

import (
	"context"
	"fmt"
	"net/http"
	"order_service/internal/cache"
	"order_service/internal/config"
	"order_service/internal/db"
	"order_service/internal/kafka"
	. "order_service/internal/logger"
	"order_service/internal/rest"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Загрузка логгер
	LoadLogger()

	// Загрузка БД
	retries := 5
	db, err := db.Load(retries)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	defer db.Close()

	// Создание кэша
	cache := cache.CreateCache(db)

	// Создание REST сервиса
	rest.CreateRestService(db)

	// Загрузка конфига
	cfg, err := config.Load()
	if err != nil {
		Logger.Fatal(err.Error())
	}

	// Запуск сервера
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

	// Запуск consumer
	kafkaConsumer := kafka.CreateConsumerWrapper(db, []string{cfg.KafkaUrl}, cache)

	go kafkaConsumer.ConsumeMessages()

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
