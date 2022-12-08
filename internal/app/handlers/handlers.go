package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/alaleks/shortener/internal/app/database/methods"
	"github.com/alaleks/shortener/internal/app/service"
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

	shortURL, err := h.AddShortenURL(userID, longURL)

	// если ошибка соот-т ErrAlreadyExists, то устанавливаем статус 409
	// в противном случае - статус 201
	switch errors.Is(err, methods.ErrAlreadyExists) {
	case true:
		writer.WriteHeader(http.StatusConflict)
	default:
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

	longURL, err := h.GetOriginalURL(uid)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.Header().Set("Location", longURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) Ping(writer http.ResponseWriter, req *http.Request) {
	err := h.PingDB()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
