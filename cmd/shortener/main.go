package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/alaleks/shortener/internal/app/serv"
)

func main() {
	sizeUID := 5
	server := serv.New(":8080", sizeUID)

	if err := serv.Run(server); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
