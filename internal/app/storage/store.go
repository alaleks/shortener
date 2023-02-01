// Package storage includes data storage implementations.
package storage

import (
	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/storage/pool"
)

// Data Structures
type (
	Store struct {
		Store Storage
		Pool  *pool.Pool
	}

	Storage interface {
		Producer
		Consumer
		User
		Worker
	}

	Statistics struct {
		ShortURL  string `json:"shorturl"`
		LongURL   string `json:"longurl"`
		CreatedAt string `json:"createdAt"`
		Usage     uint   `json:"usage"`
	}
)

// Storage interfaces
type (
	Worker interface {
		Init() error
		Close() error
		Ping() error
	}

	Producer interface {
		Add(longURL, userID string) (string, error)
		AddBatch(longURL, userID, corID string) string
		Update(uid string)
		DelUrls(userID string, shortsUID ...string) error
	}

	Consumer interface {
		GetURL(uid string) (string, error)
		Stat(uid string) (Statistics, error)
	}

	User interface {
		Create() uint
		GetUrlsUser(userID string) ([]struct {
			ShortUID string `json:"short_url"`
			LongURL  string `json:"original_url"`
		}, error)
	}
)

// InitStore initializes the store instance.
func InitStore(conf config.Configurator, logger *logger.AppLogger) *Store {
	pool := pool.Init(logger)
	go pool.Run()

	if len([]rune(conf.GetDSN())) > 1 {
		storeDB := &Store{
			Store: NewDB(conf),
			Pool:  pool,
		}

		// Initializing the database.
		err := storeDB.Store.Init()

		// Return the structure only if the error is nil
		// otherwise use file or memory.
		if err == nil {
			return storeDB
		}
	}

	storeDefault := &Store{
		Store: NewDefault(conf),
		Pool:  pool,
	}

	// Initialize file storage.
	_ = storeDefault.Store.Init()

	return storeDefault
}
