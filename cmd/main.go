package main

import (
	"log"

	"github.com/alaleks/shortener/internal/app"
)

func main() {
	srv, err := app.NewServer("8080")
	if err == nil {
		log.Fatal(app.Run(srv))
	} else {
		log.Fatal(err)
	}
}
