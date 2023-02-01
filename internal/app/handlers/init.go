package handlers

import (
	"errors"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/storage"
)

// Handler structure that includes the Storage structure.
type Handlers struct {
	Storage *storage.Store
}

// List of typical errors.
var (
	ErrEmptyURL       = errors.New("url is empty")
	ErrInternalError  = errors.New("sorry, an error has occurred, please try again")
	ErrInvalidUID     = errors.New("short url is invalid")
	ErrInvalidRequest = errors.New(`json is invalid, please check what you send. 
	Should be: {"url":"https://example.ru"}`)
	ErrUserDoesNotExist = errors.New("user did not use the service")
	ErrEmptyBatch       = errors.New("URL batching error, please check the source data")
)

// InputShorten structure for the ShortenURLAPI method containing a URL field.
type InputShorten struct {
	URL string `json:"url"`
}

// OutputShorten structure for the ShortenURLAPI method containing data for response.
type OutputShorten struct {
	Result  string `json:"result,omitempty"`
	Err     string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

// InShortenBatch structure for the ShortenURLAPI method
// containing a original URL field and Correlation ID.
type InShortenBatch struct {
	CorID       string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

// OutShortenBatch structure for the ShortenURLAPI method containing data for response.
type OutShortenBatch struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
	Err      string `json:"error,omitempty"`
}

// New - initializes handlers.
func New(conf config.Configurator, logger *logger.AppLogger) *Handlers {
	handlers := Handlers{
		Storage: storage.InitStore(conf, logger),
	}

	return &handlers
}
