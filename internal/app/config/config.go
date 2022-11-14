package config

import (
	"bytes"
	"flag"
	"os"
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
	appConf := AppConfig{
		serverAddress:   bytes.NewBuffer([]byte("localhost:8080")),
		baseURL:         bytes.NewBuffer([]byte{}),
		fileStoragePath: bytes.NewBuffer([]byte{}),
	}

	appConf.defineOptionsApp()
	appConf.checkOptions()

	return &appConf
}

func (a *AppConfig) GetServAddr() string {
	return a.serverAddress.String()
}

func (a *AppConfig) GetBaseURL() *bytes.Buffer {
	return a.baseURL
}

func (a *AppConfig) GetFileStoragePath() *bytes.Buffer {
	return a.fileStoragePath
}

func SetEnvApp(serverAddr, baseURL, pathStorage string) {
	os.Setenv("SERVER_ADDRESS", serverAddr)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("FILE_STORAGE_PATH", pathStorage)
}

func (a *AppConfig) defineOptionsApp() {
	if servAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		a.serverAddress.Reset()
		a.serverAddress.WriteString(servAddr)
	}

	if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
		a.baseURL.Reset()
		a.baseURL.WriteString(baseURL)
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		a.fileStoragePath.WriteString(fileStoragePath)
	}

	// приоритет флагов будет выше, чем установленных переменных окружения
	// если значения флагов не пустая строка, то они переопределят настройки.
	servAddr := flag.String("a", "", "SERVER_ADDRESS")
	baseURL := flag.String("b", "", "BASE_URL")
	fileStoragePath := flag.String("f", "", "FILE_STORAGE_PATH")

	flag.Parse()

	if *servAddr != "" {
		a.serverAddress.Reset()
		a.serverAddress.WriteString(*servAddr)
	}

	if *baseURL != "" {
		a.baseURL.Reset()
		a.baseURL.WriteString(*baseURL)
	}

	if *fileStoragePath != "" {
		a.fileStoragePath.WriteString(*fileStoragePath)
	}
}

func (a *AppConfig) checkOptions() {
	// проверка адреса сервера, должен быть указан порт
	if !bytes.Contains(a.serverAddress.Bytes(), []byte{58}) {
		// если порт не указан, то ставим 8080
		a.serverAddress.WriteString(":8080")
	}

	// проверка корректности base url - должен быть протокол и слэш в конце
	switch a.baseURL.Len() {
	case 0:
		a.baseURL.WriteString("http://")
		a.baseURL.Write(a.serverAddress.Bytes())
		a.baseURL.WriteString("/")
	default:
		if !bytes.HasPrefix(a.baseURL.Bytes(), []byte{104, 116, 116, 112}) {
			a.baseURL.Write(append([]byte{104, 116, 116, 112, 58, 47, 47}, a.baseURL.Next(a.baseURL.Len())...))
		}

		if !bytes.HasSuffix(a.baseURL.Bytes(), []byte{47}) {
			a.baseURL.Write([]byte{47})
		}
	}
}
