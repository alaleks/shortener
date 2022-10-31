package serv

import (
	"net/http"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/router"
)

func New(appPort string) *http.Server {
	var conf = config.New("8080")

	server := &http.Server{
		Handler: router.Create(),
		Addr:    conf.Port(),
	}

	return server
}

func Run(server *http.Server) error {
	return server.ListenAndServe()
}
