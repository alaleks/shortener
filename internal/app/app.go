package app

import (
	"net/http"

	"github.com/alaleks/shortener/internal/config"
	"github.com/alaleks/shortener/internal/router"
)

func Run() error {
	port := config.SelectAppPort("8080")

	server := &http.Server{
		Handler: router.Create(),
		Addr:    port.Get(),
	}

	return server.ListenAndServe()
}
