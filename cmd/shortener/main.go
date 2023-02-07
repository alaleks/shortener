// Application run
package main

import (
	"errors"
	"net/http"
	_ "net/http/pprof"

	"github.com/alaleks/shortener/internal/app/serv"
)

func main() {
	server := serv.New()

	// Run server for pprof
	go func() {
		server.Logger.LZ.Fatal(http.ListenAndServe(":3031", nil))
	}()

	if err := serv.Run(server); !errors.Is(err, http.ErrServerClosed) {
		server.Logger.LZ.Fatal(err)
	}
}
