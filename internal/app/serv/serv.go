package serv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/router"
	"github.com/alaleks/shortener/internal/app/serv/middleware"
)

const (
	defaultTimeout           = time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultIdleTimeout       = 15 * time.Second
	maxHeaderBytes           = 4096
)

type AppServer struct {
	server   *http.Server
	handlers *handlers.Handlers
	conf     config.Configurator
}

func New(sizeUID int) *AppServer {
	var (
		appConf    config.Configurator = config.New()
		appHandler                     = handlers.New(sizeUID, appConf)
	)

	server := &http.Server{
		Handler:           middleware.CompressHandler(router.Create(appHandler)),
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              appConf.GetServAddr(),
		TLSConfig:         nil,
		MaxHeaderBytes:    maxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	return &AppServer{server: server, handlers: appHandler, conf: appConf}
}

func Run(appServer *AppServer) error {
	go catchSignal(appServer)

	return fmt.Errorf("failed run server: %w", appServer.server.ListenAndServe())
}

func catchSignal(appServer *AppServer) {
	termSignals := make(chan os.Signal, 1)
	reloadSignals := make(chan os.Signal, 1)

	signal.Notify(termSignals, syscall.SIGINT)

	signal.Notify(reloadSignals, syscall.SIGUSR1)

	for {
		select {
		case <-termSignals:
			if appServer.conf.GetFileStoragePath().String() != "" {
				log.Fatal(appServer.handlers.DataStorage.Write(appServer.conf.GetFileStoragePath().String()))
			}

			log.Fatal(appServer.server.Shutdown(context.Background()))
		case <-reloadSignals:
			if appServer.conf.GetFileStoragePath().String() != "" {
				log.Fatal(appServer.handlers.DataStorage.Write(appServer.conf.GetFileStoragePath().String()))
			}
		}
	}
}
