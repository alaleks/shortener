package config

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

type Configurator interface {
	Recipient
	Tuner
}

type Recipient interface {
	GetServAddr() string
	GetBaseURL() *bytes.Buffer
	GetFileStoragePath() *bytes.Buffer
}

type Tuner interface {
	DefineOptionsEnv()
	DefineOptionsFlags([]string)
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
	return &AppConfig{
		serverAddress:   bytes.NewBuffer([]byte("localhost:8080")),
		baseURL:         bytes.NewBuffer([]byte{}),
		fileStoragePath: bytes.NewBuffer([]byte{}),
	}
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

func (a *AppConfig) DefineOptionsEnv() {
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

	// проверяем корректность опций
	a.checkOptions()
}

func (a *AppConfig) DefineOptionsFlags(args []string) {
	confFlags, err := parseFlags(args)

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
	// проверка адреса сервера, должен быть указан порт
	if !bytes.Contains(a.serverAddress.Bytes(), []byte{58}) {
		// если порт не указан, то ставим 8080
		a.serverAddress.WriteString(":8080")
	}

	// проверка корректности base url - должен быть протокол и слэш в конце
	switch a.baseURL.Len() {
	case 0:
		a.baseURL.Write([]byte{104, 116, 116, 112, 58, 47, 47})
		a.baseURL.Write(a.serverAddress.Bytes())
		a.baseURL.Write([]byte{47})
	default:
		if !bytes.HasPrefix(a.baseURL.Bytes(), []byte{104, 116, 116, 112}) {
			a.baseURL.Write(append([]byte{104, 116, 116, 112, 58, 47, 47}, a.baseURL.Next(a.baseURL.Len())...))
		}

		if !bytes.HasSuffix(a.baseURL.Bytes(), []byte{47}) {
			a.baseURL.Write([]byte{47})
		}
	}
}
