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
	GetSecretKey() []byte
	GetDSN() string
}

type Tuner interface {
	DefineOptionsEnv()
	DefineOptionsFlags([]string)
}

type AppConfig struct {
	serverAddress   string
	baseURL         string
	fileStoragePath string
	dsn             string
	secretKey       []byte
}

type Options struct {
	Env  bool
	Flag bool
}

type confFlags struct {
	a *string
	b *string
	f *string
	d *string
}

func New(opt Options) *AppConfig {
	appConf := AppConfig{
		serverAddress:   "localhost:8080",
		baseURL:         "http://localhost:8080/",
		fileStoragePath: "",
		dsn:             "",
		secretKey:       []byte("9EE3BF9351DFCFF24CD6DA2C4D963"),
	}

	if opt.Env {
		appConf.DefineOptionsEnv()
	}

	if opt.Flag {
		appConf.DefineOptionsFlags(os.Args)
	}

	return &appConf
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

func (a *AppConfig) GetSecretKey() []byte {
	return a.secretKey
}

func (a *AppConfig) GetDSN() string {
	return a.dsn
}

func (a *AppConfig) DefineOptionsEnv() {
	if servAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok && servAddr != "" {
		a.serverAddress = servAddr
	}

	if baseURL, ok := os.LookupEnv("BASE_URL"); ok && baseURL != "" {
		a.baseURL = baseURL
	} else {
		a.baseURL = a.serverAddress
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && fileStoragePath != "" {
		a.fileStoragePath = fileStoragePath
	}

	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok && dsn != "" {
		a.dsn = dsn
	}

	// проверяем корректность опций
	a.checkOptions()
}

func (a *AppConfig) DefineOptionsFlags(args []string) {
	confFlags, err := parseFlags(args)
	if err != nil {
		return
	}

	if *confFlags.a != "" {
		a.serverAddress = *confFlags.a
	}

	if *confFlags.b != "" {
		a.baseURL = *confFlags.b
	} else {
		a.baseURL = a.serverAddress
	}

	if *confFlags.f != "" {
		a.fileStoragePath = *confFlags.f
	}

	if *confFlags.d != "" {
		a.dsn = *confFlags.d
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
	confFlags.d = flags.String("d", "", "DATABASE_DSN")

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

	if !strings.HasPrefix(a.baseURL, "http") {
		a.baseURL = fmt.Sprintf("%s%s", httpPrefix, a.baseURL)
	}

	if !strings.HasSuffix(a.baseURL, "/") {
		a.baseURL += "/"
	}
}
