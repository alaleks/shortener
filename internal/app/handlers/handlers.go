package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"

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
	ShortenURLAPI(writer http.ResponseWriter, req *http.Request)
	GetStatAPI(writer http.ResponseWriter, req *http.Request)
}

var (
	ErrEmptyURL   = errors.New("url is empty")
	ErrWriter     = errors.New("sorry, an error has occurred, please try again")
	ErrUIDInvalid = errors.New("short url is invalid")
)

func New(sizeShortUID int) *Handlers {
	return &Handlers{DataStorage: storage.New(), SizeUID: sizeShortUID}
}

func (h *Handlers) ShortenURL(writer http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
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
