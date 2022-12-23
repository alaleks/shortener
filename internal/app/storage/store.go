package storage

import (
	"errors"
	"runtime"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/shmel1k/gop"
)

var ErrShortURLRemoved = errors.New("short URL has been removed")

type Store struct {
	Store Storage
	Pool  gop.Pool
}

type Storage interface {
	Producer
	Consumer
	User
	Worker
}

type Worker interface {
	Init() error
	Close() error
	Ping() error
}

type Producer interface {
	Add(longURL, userID string) (string, error)
	AddBatch(longURL, userID, corID string) string
	Update(uid string)
	DelUrls(userID string, shortsUID ...string) error
}

type Consumer interface {
	GetURL(uid string) (string, error)
	Stat(uid string) (Statistics, error)
}

type User interface {
	Create() uint
	GetUrlsUser(userID string) ([]URLUser, error)
}

type Statistics struct {
	ShortURL  string `json:"shorturl"`
	LongURL   string `json:"longurl"`
	CreatedAt string `json:"createdAt"`
	Usage     uint   `json:"usage"`
}

func InitStore(conf config.Configurator) *Store {
	pool := gop.NewPool(gop.Config{
		MaxWorkers:         runtime.NumCPU(),
		UnstoppableWorkers: runtime.NumCPU(),
	})

	if len([]rune(conf.GetDSN())) > 1 {
		storeDB := &Store{
			Store: NewDB(conf),
			Pool:  pool,
		}
		// инициализируем базу данных
		err := storeDB.Store.Init()

		// возвращаем структуру только если ошибка nil
		// в противном случае используем файл или память
		if err == nil {
			return storeDB
		}
	}

	storeDefault := &Store{
		Store: NewDefault(conf),
		Pool:  pool,
	}

	// инициализируем файловое хранилище
	_ = storeDefault.Store.Init()

	return storeDefault
}
