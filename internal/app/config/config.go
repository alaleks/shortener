// Package config contains configuration for application
// and functions for its settings.
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

// Configurator interface is used for application settings
// by combining interfaces Recipient and Tuner interfaces.
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
	EnableTLS() bool
}

// Tuner interface implements methods for configuring tuning.
type Tuner interface {
	DefineOptionsEnv()
	DefineOptionsFlags([]string)
}

// AppConfig struct with data for configuring the application.
type AppConfig struct {
	// servAddr is the address for server run.
	serverAddr string
	// baseURL is the URL app.
	baseURL string
	// fileStoragePath is the path for filestorage.
	fileStoragePath string
	// dsn is the database connection string.
	dsn string
	// secretKey is designed for encryption and decryption of authorization data.
	secretKey []byte
	// sizeUID sets the size of the short URL ID.
	sizeUID int
	// tls is used to enable TLS.
	tls bool
}

// The Options structure contains application configuration
// launch options with both environment variables and flags.
// If Env and Flag are set to true, the configuration
// will give priority to the flags.
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
	tls             *string
}

// New returns a pointer of struct that implements the Configurator interface.
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

// GetServAddr returns host for run server.
func (a *AppConfig) GetServAddr() string {
	return a.serverAddr
}

// GetBaseURL returns the url of the application.
func (a *AppConfig) GetBaseURL() string {
	return a.baseURL
}

// GetFileStoragePath returns path for filestorage.
func (a *AppConfig) GetFileStoragePath() string {
	return a.fileStoragePath
}

// GetSecretKey returns secret key
// for encryption and decryption of authorization data.
func (a *AppConfig) GetSecretKey() []byte {
	return a.secretKey
}

// GetDSN returns the database connection string.
func (a *AppConfig) GetDSN() string {
	return a.dsn
}

// EnableTLS returns true if TLS is enabled in config.
func (a *AppConfig) EnableTLS() bool {
	return a.tls
}

// GetSizeUID return the size of the short URL ID.
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

	if _, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		a.tls = true
	}

	// Сheck if the options are correct.
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

	if *confFlags.tls != "" {
		a.tls = true
	}

	if *confFlags.sizeUID != "" {
		i, err := strconv.Atoi(*confFlags.sizeUID)
		if err == nil && i > 3 {
			a.sizeUID = i
		}
	}

	// Сheck if the options are correct.
	a.checkOptions()
}

func parseFlags(args []string) (*confFlags, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var configFlags confFlags

	configFlags.serverAddr = flags.String("a", "", "SERVER_ADDRESS")
	configFlags.baseURL = flags.String("b", "", "BASE_URL")
	configFlags.fileStoragePath = flags.String("f", "", "FILE_STORAGE_PATH")
	configFlags.dsn = flags.String("d", "", "DATABASE_DSN")
	configFlags.sizeUID = flags.String("q", "", "SIZE_UID")
	configFlags.tls = flags.String("s", "", "ENABLE_HTTPS")

	err := flags.Parse(args[1:])
	if err != nil {
		err = fmt.Errorf("failed parse flags %w", flags.Parse(args[1:]))
	}

	return &configFlags, err
}

func (a *AppConfig) checkOptions() {
	httpPrefix := "http://"

	// Check server address, port must be specified.
	if !strings.Contains(a.serverAddr, ":") {
		// If the port is not specified, then add the default value :8080.
		a.serverAddr += ":8080"
	}

	if !strings.HasPrefix(a.baseURL, "http") {
		a.baseURL = fmt.Sprintf("%s%s", httpPrefix, a.baseURL)
	}

	if !strings.HasSuffix(a.baseURL, "/") {
		a.baseURL += "/"
	}
}
