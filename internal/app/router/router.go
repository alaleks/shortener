package router

import (
	"net/http"

	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/gorilla/mux"
)

func Create(handler handlers.Handler) http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", handler.ShortenURL).Methods("POST")
	mux.HandleFunc("/{uid}", handler.ParseShortURL).Methods("GET")
	mux.HandleFunc("/api/shorten", handler.ShortenURLAPI).Methods("POST")
	mux.HandleFunc("/api/{uid}/statistics", handler.GetStatAPI).Methods("GET")

	return mux
}
