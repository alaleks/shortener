// Package router registers application routers
package router

import (
	"net/http"

	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/gorilla/mux"
)

// Create - registers application routers
func Create(handler *handlers.Handlers) http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", handler.ShortenURL).Methods(http.MethodPost)
	mux.HandleFunc("/ping", handler.Ping).Methods(http.MethodGet)
	mux.HandleFunc("/{uid}", handler.ParseShortURL).Methods(http.MethodGet)
	mux.HandleFunc("/api/shorten", handler.ShortenURLAPI).Methods(http.MethodPost)
	mux.HandleFunc("/api/{uid}/statistics", handler.GetStatAPI).Methods(http.MethodGet)
	mux.HandleFunc("/api/user/urls", handler.GetUsersURL).Methods(http.MethodGet)
	mux.HandleFunc("/api/shorten/batch", handler.ShortenURLBatch).Methods(http.MethodPost)
	mux.HandleFunc("/api/user/urls", handler.ShortenDeletePool).Methods(http.MethodDelete)

	return mux
}
