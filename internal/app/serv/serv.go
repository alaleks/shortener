package serv

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/logger"

	"github.com/alaleks/shortener/internal/app/router"
	"github.com/alaleks/shortener/internal/app/serv/middleware"
	"github.com/alaleks/shortener/internal/app/serv/middleware/auth"
	"github.com/alaleks/shortener/internal/app/serv/middleware/compress"
)

const (
	defaultTimeout           = 2 * time.Second
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
		appConf    config.Configurator = config.New(config.Options{Env: true, Flag: true}, sizeUID)
		appHandler                     = handlers.New(appConf)
		auth                           = auth.TurnOn(appHandler.Storage, appConf.GetSecretKey())
	)

	server := &http.Server{
		Handler: middleware.New(compress.Compression, compress.Unpacking, auth.Authorization).
			Configure(router.Create(appHandler)),
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              appConf.GetServAddr(),
		TLSConfig:         nil,
		MaxHeaderBytes:    maxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	// init logger
	logger, err := logger.NewLogger()

	if err == nil {
		server.ErrorLog = log.New(logger, "", 0)
	}

	return &AppServer{
		server:   server,
		handlers: appHandler,
		conf:     appConf,
	}
}

func Run(appServer *AppServer) error {
	go catchSignal(appServer)

	return appServer.server.ListenAndServe()
}

func (appServer *AppServer) WriteLog(msg string) {
	appServer.server.ErrorLog.Fatal(msg)
}

func catchSignal(appServer *AppServer) {
	termSignals := make(chan os.Signal, 1)
	reloadSignals := make(chan os.Signal, 1)

	signal.Notify(termSignals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	signal.Notify(reloadSignals, syscall.SIGUSR1)

	for {
		select {
		case <-termSignals:
			appServer.handlers.Storage.Pool.Shutdown()

			if err := appServer.handlers.Storage.Store.Close(); err != nil {
				log.Fatal(err)
			}

			log.Fatal(appServer.server.Shutdown(context.Background()))
		case <-reloadSignals:
			if err := appServer.handlers.Storage.Store.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
