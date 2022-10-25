package main

import (
	"log"
	"net/http"

	"github.com/alaleks/shortener/internal/app"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/gorilla/mux"
)

func main() {

	mux := mux.NewRouter()
	mux.HandleFunc("/", handlers.UseShortener).Methods("POST")
	mux.HandleFunc("/{uid}", handlers.UseShortUrlRedirect).Methods("GET")
	mux.HandleFunc("/{uid}/statistic", handlers.GetStatistic).Methods("GET")
	log.Fatal(http.ListenAndServe(app.GetPort(), mux))

	// solutions without gorilla mux
	// without optimization
	/*
		http.HandleFunc("/", handlers.ShortenerWithoutMux)

		server := &http.Server{
			Addr: app.GetPort(),
		}
		log.Fatal(server.ListenAndServe())
	*/
}
