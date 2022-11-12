package serv

import (
	"net/http"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/router"
)

const (
	defaultTimeout           = time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultIdleTimeout       = 15 * time.Second
)

func New(sizeUID int) *http.Server {
	var (
		appConf    config.Configurator = config.New()
		appHandler handlers.Handler    = handlers.New(sizeUID, appConf.GetBaseURL())
	)

	server := &http.Server{
		Handler:           router.Create(appHandler),
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              appConf.GetServAddr(),
	}

	return server
}

func Run(server *http.Server) error {
	return server.ListenAndServe()
}
