package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/database/ping"
	"github.com/alaleks/shortener/internal/app/service"
	"github.com/gorilla/mux"
)

func (h *Handlers) ShortenURL(writer http.ResponseWriter, req *http.Request) {
	var userID string

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

	shortURL := h.AddShortenURL(userID, longURL)

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

	longURL, err := h.GetOriginalURL(uid)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
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

	if err := ping.Run(h.DSN); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
