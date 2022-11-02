package serv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/router"
)

func New(port string) *http.Server {
	var (
		appConf    config.Configurator = config.New(port)
		appHandler handlers.Handler    = handlers.New()
	)

	timeout, readHeaderTimeout, idleTimeout := 1, 2, 30
	server := &http.Server{
		Handler:           router.Create(appHandler),
		ReadTimeout:       time.Duration(timeout) * time.Second,
		WriteTimeout:      time.Duration(timeout) * time.Second,
		IdleTimeout:       time.Duration(idleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		Addr:              appConf.Port(),
	}

	return server
}

func Run(server *http.Server) error {
	return fmt.Errorf("server error: %w", server.ListenAndServe())
}
