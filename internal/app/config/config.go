package config

import (
	"net"
	"strings"
)

type Configurator interface {
	SelectPort(port string) error
	Port() string
}

type AppConfig struct {
	port string
}

func (a *AppConfig) SelectPort(port string) error {
	// checking correct port val
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	// check if the port is available
	ln, err := net.Listen("tcp", port)

	if ln != nil {
		defer ln.Close()
	}

	if err == nil {
		a.port = port
	}

	return err

}

func (a *AppConfig) Port() string {
	return a.port
}

func New() Configurator {
	return &AppConfig{}
}
