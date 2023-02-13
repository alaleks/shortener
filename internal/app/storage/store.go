// Package storage includes data storage implementations.
package storage

import (
	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/storage/pool"
)

// Data types and interfaces.
type (
	// Store represents the application's storage structure and
	// includes the Storage interface and a pointer to Pool.
	Store struct {
		Store Storage
		Pool  *pool.Pool
	}

	// Statistics represents a data model for getting statistics
	// for a specific short link.
	Statistics struct {
		ShortURL  string `json:"shorturl"`
		LongURL   string `json:"longurl"`
		CreatedAt string `json:"createdAt"`
		Usage     uint   `json:"usage"`
	}

	// Storage interface is construct to create an application's storage
	Storage interface {
		Producer
		Consumer
		User
		Worker
	}

	// Worker interface is used to initialize, ping and close application's storage.
	Worker interface {
		Init() error
		Close() error
		Ping() error
	}

	// Producer interface is used adding, updating and deleting data from application's storage.
	Producer interface {
		Add(longURL, userID string) (string, error)
		AddBatch(longURL, userID, corID string) string
		Update(uid string)
		DelUrls(userID string, shortsUID ...string) error
	}

	// Consumer interface is used gettings data from application's storage.
	Consumer interface {
		GetURL(uid string) (string, error)
		Stat(uid string) (Statistics, error)
	}

	// User interface is used to get user data from application's storage or create new user.
	User interface {
		Create() uint
		GetUrlsUser(userID string) ([]struct {
			ShortUID string `json:"short_url"`
			LongURL  string `json:"original_url"`
		}, error)
	}
)

// InitStore performs initializing the store instance.
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
