package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	DB    *gorm.DB
	sqlDB *sql.DB
}

func NewConnection() (*Connection, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getenv("DB_USER", "proy2"),
		getenv("DB_PASSWORD", "secret"),
		getenv("DB_HOST", "db"),
		getenv("DB_PORT", "5432"),
		getenv("DB_NAME", "proyecto2"),
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	var lastErr error
	for range 15 {
		if err := sqlDB.Ping(); err == nil {
			return &Connection{DB: gormDB, sqlDB: sqlDB}, nil
		} else {
			lastErr = err
		}
		time.Sleep(time.Second)
	}

	sqlDB.Close()
	return nil, lastErr
}

func (c *Connection) Close() error {
	return c.sqlDB.Close()
}

func (c *Connection) Ping(ctx context.Context) error {
	return c.sqlDB.PingContext(ctx)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
