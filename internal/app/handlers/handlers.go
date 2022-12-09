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

func (h *Handlers) ShortenURL(writer http.ResponseWriter, req *http.Request) {
	var userID string

	body, err := io.ReadAll(req.Body)

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

	shortURL, err := h.Storage.Store.Add(longURL, userID)

	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			writer.WriteHeader(http.StatusConflict)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}
	} else {
		writer.WriteHeader(http.StatusCreated)
	}

	if _, err := writer.Write([]byte(shortURL)); err != nil {
		http.Error(writer, ErrInternalError.Error(), http.StatusBadRequest)

		return
	}
}

func (h *Handlers) ParseShortURL(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]

	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	longURL, err := h.Storage.Store.GetURL(uid)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	h.Storage.Store.Update(uid)

	writer.Header().Set("Location", longURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) Ping(writer http.ResponseWriter, req *http.Request) {
	err := h.Storage.Store.Ping()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
