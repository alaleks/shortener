package serv

import (
	"net/http"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/router"
)

func New(appPort string) (*http.Server, error) {
	var conf = config.New()

	err := conf.SelectPort("8080")

	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Handler: router.Create(),
		Addr:    conf.Port(),
	}

	return server, err
}

func Run(server *http.Server) error {
	return server.ListenAndServe()
}
