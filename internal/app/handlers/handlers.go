package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/database"
	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

type Handlers struct {
	baseURL     string
	DSN         string
	DataStorage storage.Storage
	Users       storage.Users
	SizeUID     int
}

var (
	ErrEmptyURL   = errors.New("url is empty")
	ErrWriter     = errors.New("sorry, an error has occurred, please try again")
	ErrUIDInvalid = errors.New("short url is invalid")
)

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

func (h *Handlers) ShortenURL(writer http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if req.Body != nil {
		defer req.Body.Close()
	}

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	longURL := string(bytes.TrimSpace(body))

	if longURL == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	err = service.IsURL(longURL)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.WriteHeader(http.StatusCreated)

	// формируем короткую ссылку
	shortURL := h.baseURL
	uid := h.DataStorage.Add(longURL, h.SizeUID)

	if req.URL.User != nil {
		h.Users.AddShortUID(req.URL.User.Username(), uid)
	}

	shortURL += uid

	if _, err := writer.Write([]byte(shortURL)); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) ParseShortURL(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]

	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	longURL, ok := h.DataStorage.GetURL(uid)
	if !ok {
		http.Error(writer, ErrUIDInvalid.Error(), http.StatusBadRequest)

		return
	}

	h.DataStorage.Update(uid)
	writer.Header().Set("Location", longURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) Ping(writer http.ResponseWriter, req *http.Request) {
	// host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432
	if h.DSN == "" {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	err := database.CheckConnect(h.DSN)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
