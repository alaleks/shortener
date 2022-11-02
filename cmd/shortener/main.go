package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/alaleks/shortener/internal/app/serv"
)

func main() {
	server := serv.New(":8080")

	if err := serv.Run(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
