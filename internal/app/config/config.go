// The config package implements the ability to configure an application.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultSizeUID = 5
)

// Configurator interface including Recipient and Tuner interfaces.
type Configurator interface {
	Recipient
	Tuner
}

// Recipient interface implements methods for getting settings parameters.
type Recipient interface {
	GetServAddr() string
	GetBaseURL() string
	GetFileStoragePath() string
	GetSecretKey() []byte
	GetDSN() string
	GetSizeUID() int
}

// Tuner interface implements methods for configuring tuning.
type Tuner interface {
	DefineOptionsEnv()
	DefineOptionsFlags([]string)
}

// AppConfig structure that contains the app's settings settings.
type AppConfig struct {
	serverAddr      string
	baseURL         string
	fileStoragePath string
	dsn             string
	secretKey       []byte
	sizeUID         int
}

// The Options structure contains application configuration
// launch options with both environment variables and flags.
type Options struct {
	Env  bool
	Flag bool
}

type confFlags struct {
	serverAddr      *string
	baseURL         *string
	fileStoragePath *string
	dsn             *string
	sizeUID         *string
}

// New implements initialize settings.
func New(opt Options) *AppConfig {
	appConf := AppConfig{
		serverAddr:      "localhost:8080",
		baseURL:         "http://localhost:8080/",
		fileStoragePath: "",
		dsn:             "",
		secretKey:       []byte("9EE3BF9351DFCFF24CD6DA2C4D963"),
		sizeUID:         defaultSizeUID,
	}

	if opt.Env {
		appConf.DefineOptionsEnv()
	}

	if opt.Flag {
		appConf.DefineOptionsFlags(os.Args)
	}

	return &appConf
}

// GetServAddr return host for run server.
func (a *AppConfig) GetServAddr() string {
	return a.serverAddr
}

// GetBaseURL returns the url of the application.
func (a *AppConfig) GetBaseURL() string {
	return a.baseURL
}

// GetFileStoragePath return Path for Filestorage.
func (a *AppConfig) GetFileStoragePath() string {
	return a.fileStoragePath
}

// GetSecretKey return Secret Key for encrypt/decrypt.
func (a *AppConfig) GetSecretKey() []byte {
	return a.secretKey
}

// GetDSN return Data source name for Database.
func (a *AppConfig) GetDSN() string {
	return a.dsn
}

// GetSizeUID return size for create UID.
func (a *AppConfig) GetSizeUID() int {
	return a.sizeUID
}

// DefineOptionsEnv implements application configuration using environment variables.
func (a *AppConfig) DefineOptionsEnv() {
	if servAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok && servAddr != "" {
		a.serverAddr = servAddr
	}

	if baseURL, ok := os.LookupEnv("BASE_URL"); ok && baseURL != "" {
		a.baseURL = baseURL
	} else {
		a.baseURL = a.serverAddr
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && fileStoragePath != "" {
		a.fileStoragePath = fileStoragePath
	}

	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok && dsn != "" {
		a.dsn = dsn
	}

	if sizeUID, ok := os.LookupEnv("SIZE_UID"); ok && sizeUID != "" {
		i, err := strconv.Atoi(sizeUID)
		if err == nil && i > 3 {
			a.sizeUID = i
		}
	}

	// проверяем корректность опций
	a.checkOptions()
}

// DefineOptionsFlags implements application configuration using flags.
func (a *AppConfig) DefineOptionsFlags(args []string) {
	confFlags, err := parseFlags(args)
	if err != nil {
		return
	}

	if *confFlags.serverAddr != "" {
		a.serverAddr = *confFlags.serverAddr
	}

	if *confFlags.baseURL != "" {
		a.baseURL = *confFlags.baseURL
	} else {
		a.baseURL = a.serverAddr
	}

	if *confFlags.fileStoragePath != "" {
		a.fileStoragePath = *confFlags.fileStoragePath
	}

	if *confFlags.dsn != "" {
		a.dsn = *confFlags.dsn
	}

	if *confFlags.sizeUID != "" {
		i, err := strconv.Atoi(*confFlags.sizeUID)
		if err == nil && i > 3 {
			a.sizeUID = i
		}
	}

	// проверяем корректность опций
	a.checkOptions()
}

func parseFlags(args []string) (*confFlags, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var confFlags confFlags

	confFlags.serverAddr = flags.String("a", "", "SERVER_ADDRESS")
	confFlags.baseURL = flags.String("b", "", "BASE_URL")
	confFlags.fileStoragePath = flags.String("f", "", "FILE_STORAGE_PATH")
	confFlags.dsn = flags.String("d", "", "DATABASE_DSN")
	confFlags.sizeUID = flags.String("s", "", "SIZE_UID")

	err := flags.Parse(args[1:])
	if err != nil {
		err = fmt.Errorf("failed parse flags %w", flags.Parse(args[1:]))
	}

	return &confFlags, err
}

func (a *AppConfig) checkOptions() {
	httpPrefix := "http://"

	// проверка адреса сервера, должен быть указан порт
	if !strings.Contains(a.serverAddr, ":") {
		// если порт не указан, то добавляем 8080
		a.serverAddr += ":8080"
	}

	if !strings.HasPrefix(a.baseURL, "http") {
		a.baseURL = fmt.Sprintf("%s%s", httpPrefix, a.baseURL)
	}

	if !strings.HasSuffix(a.baseURL, "/") {
		a.baseURL += "/"
	}
}
