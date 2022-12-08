package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MaxIdleConns = 5
	MaxOpenConns = 50
	MaxLifetime  = time.Hour
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	sqlDB, err := db.DB()

	if err == nil {
		sqlDB.SetMaxIdleConns(MaxIdleConns)
		sqlDB.SetMaxOpenConns(MaxOpenConns)
		sqlDB.SetConnMaxLifetime(MaxLifetime)
	}

	return db, nil
}

func Ping(sqlDB *sql.DB) error {
	return sqlDB.Ping()
}
