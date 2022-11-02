package config

import (
	"strings"
)

type Configurator interface {
	Port() string
}

type AppConfig struct {
	port string
}

func (a *AppConfig) Port() string {
	return a.port
}

func New(port string) *AppConfig {
	// checking correct port val
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return &AppConfig{port: port}
}
