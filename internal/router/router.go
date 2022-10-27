package router

import (
	"net/http"

	"github.com/alaleks/shortener/internal/handlers"
	"github.com/gorilla/mux"
)

func Create() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", handlers.UseShortner).Methods("POST")
	mux.HandleFunc("/{uid}", handlers.ParseShortUrl).Methods("GET")
	mux.HandleFunc("/{uid}/statistic", handlers.GiveAwayStatistic).Methods("GET")
	return mux
}
