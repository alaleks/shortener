package config

import (
	"bytes"
	"os"
	"strings"
)

type Configurator interface {
	GetServAddr() string
	GetBaseURL() bytes.Buffer
}

type AppConfig struct {
	serverAddress bytes.Buffer
	baseURL       bytes.Buffer
}

func (a *AppConfig) GetServAddr() string {
	return a.serverAddress.String()
}

func (a *AppConfig) GetBaseURL() bytes.Buffer {
	return a.baseURL
}

func New() *AppConfig {
	var (
		appconf = AppConfig{
			serverAddress: *bytes.NewBuffer([]byte("localhost:8080")),
			baseURL:       *bytes.NewBuffer([]byte("localhost:8080")),
		}

		servAddr   string
		baseURL    string
		okServAddr bool
		okBaseURL  bool
	)

	servAddr, okServAddr = os.LookupEnv("SERVER_ADDRESS")
	baseURL, okBaseURL = os.LookupEnv("BASE_URL")

	if okServAddr {
		appconf.serverAddress.Reset()

		if !strings.Contains(servAddr, ":") {
			servAddr += ":8080"
		}

		appconf.serverAddress.WriteString(servAddr)
	}

	if okBaseURL {
		appconf.baseURL.Reset()

		if !strings.Contains(baseURL, "://") {
			appconf.baseURL.WriteString("http://")
		}

		appconf.baseURL.WriteString(baseURL)

		if !bytes.HasSuffix(appconf.baseURL.Bytes(), []byte{47}) {
			appconf.baseURL.WriteString("/")
		}
	}

	return &appconf
}

func SetEnvApp(serverAddr, baseURL string) {
	os.Setenv("SERVER_ADDRESS", serverAddr)
	os.Setenv("BASE_URL", baseURL)
}
