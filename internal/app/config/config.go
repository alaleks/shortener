package config

import (
	"bytes"
	"flag"
	"fmt"
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

type Options struct {
	Env  bool
	Flag bool
}

type confFlags struct {
	a *string
	b *string
	f *string
}

func New(opt *Options) *AppConfig {
	appConf := AppConfig{
		serverAddress:   bytes.NewBuffer([]byte("localhost:8080")),
		baseURL:         bytes.NewBuffer([]byte{}),
		fileStoragePath: bytes.NewBuffer([]byte{}),
	}

	if opt != nil {
		if opt.Env {
			appConf.defineOptionsEnv()
		}

		if opt.Flag {
			// приоритет флагов будет выше, чем установленных переменных окружения
			// если значения флагов не пустая строка, то они переопределят настройки.
			appConf.defineOptionsFlags()
		}
	}

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

func (a *AppConfig) defineOptionsEnv() {
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
}

func (a *AppConfig) defineOptionsFlags() {
	confFlags, err := parseFlags()

	if err == nil {
		if *confFlags.a != "" {
			a.serverAddress.Reset()
			a.serverAddress.WriteString(*confFlags.a)
		}

		if *confFlags.b != "" {
			a.baseURL.Reset()
			a.baseURL.WriteString(*confFlags.b)
		}

		if *confFlags.f != "" {
			a.fileStoragePath.WriteString(*confFlags.f)
		}
	}
}

func parseFlags() (*confFlags, error) {
	args := os.Args

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var confFlags confFlags

	switch strings.HasSuffix(args[0], ".test") {
	case false:
		confFlags.a = flags.String("a", "", "SERVER_ADDRESS")
		confFlags.b = flags.String("b", "", "BASE_URL")
		confFlags.f = flags.String("f", "", "FILE_STORAGE_PATH")
	case true:
		confFlags.a = flags.String("a", "", "SERVER_ADDRESS")
		confFlags.b = flags.String("b", "", "BASE_URL")
		confFlags.f = flags.String("f", "", "FILE_STORAGE_PATH")
		flags.String("test.paniconexit0", "", "TEST_PANIC_FLAG")
	}

	confFlags.a = flags.String("a", "", "SERVER_ADDRESS")
	confFlags.b = flags.String("b", "", "BASE_URL")
	confFlags.f = flags.String("f", "", "FILE_STORAGE_PATH")
	flags.String("test.paniconexit0", "", "TEST_PANIC_FLAG")

	if strings.HasSuffix(args[0], ".test") {
		flags.String("test.paniconexit0", "", "TEST_PANIC_FLAG")
	}

	err := flags.Parse(args[1:])
	if err != nil {
		err = fmt.Errorf("failed parse flags %w", flags.Parse(args[1:]))
	}

	return &confFlags, err
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
