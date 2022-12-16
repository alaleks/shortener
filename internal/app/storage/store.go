package storage

import (
	"errors"

	"github.com/alaleks/shortener/internal/app/config"
)

var ErrShortURLDeleted = errors.New("short URL has been removed")

type Store struct {
	Store Storage
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
	DelUrls(userID string, shortsUID ...string)
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
	if len([]rune(conf.GetDSN())) > 1 {
		storeDB := &Store{Store: NewDB(conf)}
		// инициализируем базу данных
		err := storeDB.Store.Init()

		// возвращаем структуру только если ошибка nil
		// в противном случае используем файл или память
		if err == nil {
			return storeDB
		}
	}

	storeDefault := &Store{Store: NewDefault(conf)}
	// инициализируем файловое хранилище
	_ = storeDefault.Store.Init()

	return storeDefault
}
