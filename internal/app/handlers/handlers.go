package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

type Handlers struct {
	DataStorage storage.Storage
	SizeUID     int
}

type Handler interface {
	ShortenURL(writer http.ResponseWriter, req *http.Request)
	ParseShortURL(writer http.ResponseWriter, req *http.Request)
	GetStat(writer http.ResponseWriter, req *http.Request)
}

var (
	ErrEmptyURL   = errors.New("url is empty")
	ErrWriter     = errors.New("sorry, an error has occurred, please try again")
	ErrUIDInvalid = errors.New("short url is invalid")
)

type Statistics struct {
	ShortURL  string `json:"shorturl"`
	LongURL   string `json:"longurl"`
	Usage     uint   `json:"usage"`
	CreatedAt string `json:"createdAt"`
}

func New(size int) *Handlers {
	return &Handlers{DataStorage: storage.New(), SizeUID: size}
}

func (h *Handlers) ShortenURL(writer http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	longURL := strings.TrimSpace(string(body))

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
	var shortURL bytes.Buffer

	shortURL.WriteString(func() string {
		if req.TLS != nil {
			return "https://"
		}

		return "http://"
	}())

	shortURL.WriteString(req.Host + "/")

	uid := h.DataStorage.Add(longURL, h.SizeUID)
	shortURL.WriteString(uid)

	if _, err := writer.Write(shortURL.Bytes()); err != nil {
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

func (h *Handlers) GetStat(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]
	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	host := "http://" + req.Host + "/"

	if req.TLS != nil {
		host = "https://" + req.Host + "/"
	}

	longURL, counterStat, createdAt := h.DataStorage.Stat(uid)

	if longURL == "" {
		http.Error(writer, ErrUIDInvalid.Error(), http.StatusBadRequest)

		return
	}

	dataForRes := Statistics{ShortURL: host + uid, LongURL: longURL, Usage: counterStat, CreatedAt: createdAt}
	toJSON, err := json.Marshal(dataForRes)

	if err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err := writer.Write(toJSON); err != nil {
		http.Error(writer, ErrWriter.Error(), http.StatusBadRequest)

		return
	}
}
