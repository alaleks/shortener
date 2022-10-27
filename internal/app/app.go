package app

import (
	"net"
	"net/http"

	"github.com/alaleks/shortener/internal/config"
	"github.com/alaleks/shortener/internal/router"
)

func NewServer(appPort string) (*http.Server, error) {
	port := config.SelectAppPort(appPort)

	ln, err := net.Listen("tcp", port.Get())

	if err != nil {
		return nil, err
	}

	defer ln.Close()

	server := &http.Server{
		Handler: router.Create(),
		Addr:    port.Get(),
	}

	return server, err
}

func Run(server *http.Server) error {
	return server.ListenAndServe()
}
