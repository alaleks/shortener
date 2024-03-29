// Package handlers implements application route handlers
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

// ShortenURL implements URL shortening.
//
// POST /, text: "http://github.com/alaleks/shortener".
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

	shortURL, err := h.Storage.St.Add(longURL, userID)

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

// ParseShortURL takes a short URL and redirects at the original URL.
//
// GET /{uid}
func (h *Handlers) ParseShortURL(writer http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]

	if uid == "" {
		http.Error(writer, ErrEmptyURL.Error(), http.StatusBadRequest)

		return
	}

	longURL, err := h.Storage.St.GetURL(uid)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, storage.ErrShortURLRemoved) {
			status = http.StatusGone
		}

		http.Error(writer, err.Error(), status)

		return
	}

	h.Storage.St.Update(uid)

	writer.Header().Set("Location", longURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

// Ping test the application.
//
// GET /ping
func (h *Handlers) Ping(writer http.ResponseWriter, req *http.Request) {
	if err := h.Storage.St.Ping(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
