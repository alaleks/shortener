package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	dbInstance, err := db.DB()
	if err != nil {
		return err
	}

	return dbInstance.Close()
}

func CheckConnect(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return err
	}

	return nil
}
