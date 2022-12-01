package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/database/methods"
	"github.com/alaleks/shortener/internal/app/database/ping"
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
	var (
		shortUID string
		userID   string
	)

	body, err := io.ReadAll(req.Body)

	if req.Body != nil {
		defer req.Body.Close()
	}

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	longURL := string(bytes.TrimSpace(body))
	err = service.IsURL(longURL)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if req.URL.User != nil {
		userID = req.URL.User.Username()
	}

	writer.WriteHeader(http.StatusCreated)

	// формируем короткую ссылку
	shortURL := h.baseURL

	if h.DSN != "" {
		d := methods.NewDB(h.DSN)

		if d.DB != nil {
			shortUID = service.GenUID(h.SizeUID)
			d.AddURL(userID, shortUID, longURL)

			defer d.Close()
		}
	} else {
		shortUID = h.DataStorage.Add(longURL, h.SizeUID)
		h.Users.AddShortUID(userID, shortUID)
	}

	shortURL += shortUID

	if _, err := writer.Write([]byte(shortURL)); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) ParseShortURL(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]

	var longURL string
	var ok bool

	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	switch h.DSN {
	case "":
		longURL, ok = h.DataStorage.GetURL(uid)

		if !ok {
			http.Error(writer, ErrUIDInvalid.Error(), http.StatusBadRequest)

			return
		}

		h.DataStorage.Update(uid)
	default:
		d := methods.NewDB(h.DSN)

		if d.DB != nil {
			longURL = d.GetOriginalURL(uid)

			if longURL == "" {
				http.Error(writer, ErrUIDInvalid.Error(), http.StatusBadRequest)

				return
			}

			d.UpdateStat(uid)

			defer d.Close()
		}
	}

	writer.Header().Set("Location", longURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) Ping(writer http.ResponseWriter, req *http.Request) {
	// host=localhost user=shortener password=3BJ2zWGPbQps dbname=shortener port=5432
	if h.DSN == "" {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	err := ping.Run(h.DSN)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
