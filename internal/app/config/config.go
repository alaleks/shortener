package config

import (
	"bytes"
	"os"
	"strings"
)

type Configurator interface {
	GetServAddr() string
	GetBaseURL() *bytes.Buffer
	GetFileStoragePath() *bytes.Buffer
}

type AppConfig struct {
	serverAddress   *bytes.Buffer
	baseURL         *bytes.Buffer
	fileStoragePath *bytes.Buffer
}

func New() *AppConfig {
	return &AppConfig{
		serverAddress:   bytes.NewBuffer([]byte("localhost:8080")),
		baseURL:         bytes.NewBuffer([]byte("http://localhost:8080/")),
		fileStoragePath: bytes.NewBuffer([]byte{}),
	}
}

func (a *AppConfig) GetServAddr() string {
	servAddr, ok := os.LookupEnv("SERVER_ADDRESS")

	if !ok {
		return a.serverAddress.String()
	}

	a.serverAddress.Reset()

	if !strings.Contains(servAddr, ":") {
		servAddr += ":8080"
	}

	a.serverAddress.WriteString(servAddr)

	return a.serverAddress.String()
}

func (a *AppConfig) GetBaseURL() *bytes.Buffer {
	baseURL, ok := os.LookupEnv("BASE_URL")

	if !ok {
		return a.baseURL
	}

	a.baseURL.Reset()

	if !strings.Contains(baseURL, "://") {
		a.baseURL.WriteString("http://")
	}

	a.baseURL.WriteString(baseURL)

	if !bytes.HasSuffix(a.baseURL.Bytes(), []byte{47}) {
		a.baseURL.WriteString("/")
	}

	return a.baseURL
}

func (a *AppConfig) GetFileStoragePath() *bytes.Buffer {
	fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")

	if ok {
		a.fileStoragePath.WriteString(fileStoragePath)
	}

	return a.fileStoragePath
}

func SetEnvApp(serverAddr, baseURL, pathStorage string) {
	os.Setenv("SERVER_ADDRESS", serverAddr)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("FILE_STORAGE_PATH", pathStorage)
}
