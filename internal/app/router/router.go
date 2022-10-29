package router

import (
	"net/http"

	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/gorilla/mux"
)

func Create() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", handlers.ShortenURL).Methods("POST")
	mux.HandleFunc("/{uid}", handlers.ParseShortURL).Methods("GET")
	mux.HandleFunc("/{uid}/statistic", handlers.GetStat).Methods("GET")
	return mux
}
