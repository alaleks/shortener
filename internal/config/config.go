package config

import (
	"net"
	"strings"
)

type Configer interface {
	SelectPort(port string) error
	Port() string
}

type AppConfig struct {
	port string
}

func (a *AppConfig) SelectPort(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
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

func New() Configer {
	return &AppConfig{}
}
