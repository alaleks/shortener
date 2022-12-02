package handlers

import (
	"errors"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/storage"
)

type Handlers struct {
	baseURL     string
	DSN         string
	DataStorage storage.Storage
	Users       storage.Users
	SizeUID     int
}

var (
	ErrEmptyURL    = errors.New("url is empty")
	ErrWriter      = errors.New("sorry, an error has occurred, please try again")
	ErrUIDInvalid  = errors.New("short url is invalid")
	ErrInvalidJSON = errors.New(`json is invalid, please check what you send. 
	Should be: {"url":"https://example.ru"}`)
	ErrGetUrlsUser = errors.New("user did not use the service")
	ErrEmptyBatch  = errors.New("URL batching error, please check the source data")
)

type Statistics struct {
	ShortURL  string `json:"shorturl"`
	LongURL   string `json:"longurl"`
	Usage     uint   `json:"usage"`
	CreatedAt string `json:"createdAt"`
}

type InputShorten struct {
	URL string `json:"url"`
}

type OutputShorten struct {
	Result  string `json:"result,omitempty"`
	Err     string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

type InShortenBatch struct {
	CorID       string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type OutShortenBatch struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
	Err      string `json:"error,omitempty"`
}

func New(sizeShortUID int, conf config.Configurator) *Handlers {
	handlers := Handlers{
		DataStorage: storage.New(),
		SizeUID:     sizeShortUID,
		baseURL:     conf.GetBaseURL(),
		Users:       storage.NewUsers(),
		DSN:         conf.GetDSN(),
	}

	if conf.GetFileStoragePath() != "" {
		err := handlers.DataStorage.Read(conf.GetFileStoragePath())
		if err != nil {
			return &handlers
		}
	}

	return &handlers
}
