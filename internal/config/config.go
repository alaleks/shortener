package config

import (
	"strings"
)

type AppPort string

func SelectAppPort(port string) *AppPort {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	return (*AppPort)(&port)
}

func (a *AppPort) Get() string {
	return string(*a)
}
