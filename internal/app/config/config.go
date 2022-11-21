package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Configurator interface {
	Recipient
	Tuner
}

type Recipient interface {
	GetServAddr() string
	GetBaseURL() string
	GetFileStoragePath() string
}

type Tuner interface {
	DefineOptionsEnv()
	DefineOptionsFlags([]string)
}

type AppConfig struct {
	serverAddress   string
	baseURL         string
	fileStoragePath string
}

type Options struct {
	Env  bool
	Flag bool
}

type confFlags struct {
	a *string
	b *string
	f *string
}

func New(opt Options) *AppConfig {
	return &AppConfig{
		serverAddress:   "localhost:8080",
		baseURL:         "",
		fileStoragePath: "",
	}
}

func (a *AppConfig) GetServAddr() string {
	return a.serverAddress
}

func (a *AppConfig) GetBaseURL() string {
	return a.baseURL
}

func (a *AppConfig) GetFileStoragePath() string {
	return a.fileStoragePath
}

func (a *AppConfig) DefineOptionsEnv() {
	if servAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		a.serverAddress = servAddr
	}

	if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
		a.baseURL = baseURL
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		a.fileStoragePath = fileStoragePath
	}

	// проверяем корректность опций
	a.checkOptions()
}

func (a *AppConfig) DefineOptionsFlags(args []string) {
	confFlags, err := parseFlags(args)

	if err == nil {
		if *confFlags.a != "" {
			a.serverAddress = *confFlags.a
		}

		if *confFlags.b != "" {
			a.baseURL = *confFlags.b
		}

		if *confFlags.f != "" {
			a.fileStoragePath = *confFlags.f
		}
	}

	// проверяем корректность опций
	a.checkOptions()
}

func parseFlags(args []string) (*confFlags, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var confFlags confFlags

	confFlags.a = flags.String("a", "", "SERVER_ADDRESS")
	confFlags.b = flags.String("b", "", "BASE_URL")
	confFlags.f = flags.String("f", "", "FILE_STORAGE_PATH")

	err := flags.Parse(args[1:])
	if err != nil {
		err = fmt.Errorf("failed parse flags %w", flags.Parse(args[1:]))
	}

	return &confFlags, err
}

func (a *AppConfig) checkOptions() {
	httpPrefix := "http://"

	// проверка адреса сервера, должен быть указан порт
	if !strings.Contains(a.serverAddress, ":") {
		// если порт не указан, то добавляем 8080
		a.serverAddress += ":8080"
	}

	// проверка корректности base url - должен быть протокол и слэш в конце
	switch len(a.baseURL) {
	case 0:
		a.baseURL = httpPrefix + a.serverAddress + "/"
	default:
		if !strings.HasPrefix(a.baseURL, "http") {
			uri := a.baseURL
			a.baseURL = httpPrefix + uri
		}

		if !strings.HasSuffix(a.baseURL, "/") {
			a.baseURL += "/"
		}
	}
}
