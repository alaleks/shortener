package handlers

import (
	"errors"
	"fmt"

	"github.com/alaleks/shortener/internal/app/database"
	"github.com/alaleks/shortener/internal/app/database/methods"
	_ "github.com/lib/pq"
)

var ErrDatabaseConnection = errors.New("database connection not established when application starts")

func (h *Handlers) ConnectDB() error {
	db, err := database.Connect(h.DSN)

	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	h.DB = &methods.Database{SDB: db}
	h.checkDB = true

	return nil
}

func (h *Handlers) PingDB() error {
	if !h.checkDB {
		return ErrDatabaseConnection
	}

	sqlDB, err := h.DB.SDB.DB()

	if err != nil {
		return h.ConnectDB()
	}

	err = sqlDB.Ping()

	if err != nil {
		return h.ConnectDB()
	}

	return nil
}

func (h *Handlers) CloseDB() error {
	if !h.checkDB {
		return nil
	}

	return h.DB.Close()
}
