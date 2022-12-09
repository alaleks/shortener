package handlers

import (
	"errors"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/storage"
)

type Handlers struct {
	Storage *storage.Store
}

var (
	ErrEmptyURL       = errors.New("url is empty")
	ErrInternalError  = errors.New("sorry, an error has occurred, please try again")
	ErrInvalidUID     = errors.New("short url is invalid")
	ErrInvalidRequest = errors.New(`json is invalid, please check what you send. 
	Should be: {"url":"https://example.ru"}`)
	ErrUserDoesNotExist = errors.New("user did not use the service")
	ErrEmptyBatch       = errors.New("URL batching error, please check the source data")
)

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

func New(conf config.Configurator) *Handlers {
	handlers := Handlers{
		Storage: storage.InitStore(conf),
	}

	return &handlers
}
