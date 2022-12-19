package main

import (
	"errors"
	"net/http"

	"github.com/alaleks/shortener/internal/app/serv"
)

func main() {
	sizeUID := 5
	server := serv.New(sizeUID)
	fatalLog := server.AppLogger.Fatal()

	if err := serv.Run(server); !errors.Is(err, http.ErrServerClosed) {
		fatalLog.Println(err)
	}
}
