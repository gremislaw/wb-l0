package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"order_service/internal/config"
	. "order_service/internal/logger"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// Загрузка БД
func Load(retries int) (*sql.DB, error) {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		Logger.Warn(err.Error())
		return nil, err
	}

	// Преобразование конфигурационных данных в DSN
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.POSTGRES_HOST,
		cfg.POSTGRES_PORT, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB)

	db, err := Connect(connStr, retries)
	if err != nil {
		Logger.Warn(err.Error())
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		Logger.Warn(err.Error())
		return nil, err
	}

	Logger.Info("Database successfully loaded", zap.String("DSN", connStr))
	return db, nil
}

// Подключение к БД
func Connect(dsn string, retries int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 1; i <= retries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			break
		}
		errMsg := "DB connection error"
		Logger.Warn(errMsg, zap.String("DSN", dsn))
		Logger.Info("retrying to load DB...", zap.Int("retry", i))
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		errMsg := "DB connection error"
		Logger.Warn(errMsg, zap.String("DSN", dsn))
		return nil, errors.New(errMsg)
	}
	return db, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Миграция БД
func Migrate(DB *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		errMsg := "Failed to set postgres dialect for goose: "
		Logger.Warn(errMsg)
		return errors.New(errMsg + err.Error())
	}
	if err := goose.Up(DB, "migrations"); err != nil {
		errMsg := "Failed to migrate database: "
		Logger.Warn(errMsg)
		return errors.New(errMsg + err.Error())
	}

	return nil
}
