package database

import (
	"database/sql"
	"ecommerce/util"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func InitDB(config util.Configuration, logger *zap.Logger) (*sql.DB, error) {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Host)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Panic("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Panic("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection initialized successfully")
	return db, nil
}
