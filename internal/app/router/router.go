package router

import (
	"net/http"

	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

func Create() http.Handler {
	DataStorage := storage.New()

	mux := mux.NewRouter()
	mux.HandleFunc("/", handlers.ShortenURL(DataStorage)).Methods("POST")
	mux.HandleFunc("/{uid}", handlers.ParseShortURL(DataStorage)).Methods("GET")
	mux.HandleFunc("/{uid}/statistic", handlers.GetStat(DataStorage)).Methods("GET")
	return mux
}
