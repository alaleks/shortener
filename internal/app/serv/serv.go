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
	applogger "github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/router"
	"github.com/alaleks/shortener/internal/app/serv/middleware"
	"github.com/alaleks/shortener/internal/app/serv/middleware/auth"
	"github.com/alaleks/shortener/internal/app/serv/middleware/compress"
	muxhandlers "github.com/gorilla/handlers"
)

const (
	defaultTimeout           = 2 * time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultIdleTimeout       = 15 * time.Second
	maxHeaderBytes           = 4096
	timeOutClosedServer      = 2 * time.Second
)

type AppServer struct {
	server    *http.Server
	handlers  *handlers.Handlers
	conf      config.Configurator
	AppLogger *applogger.LogDir
}

func New(sizeUID int) *AppServer {
	var (
		appConf    config.Configurator = config.New(config.Options{Env: true, Flag: true}, sizeUID)
		applogger                      = applogger.New()
		appHandler                     = handlers.New(appConf)
		auth                           = auth.TurnOn(appHandler.Storage, appConf.GetSecretKey())
	)

	server := &http.Server{
		Handler: muxhandlers.RecoveryHandler()(middleware.New(compress.Compression, compress.Unpacking, auth.Authorization).
			Configure(router.Create(appHandler))),
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              appConf.GetServAddr(),
		TLSConfig:         nil,
		MaxHeaderBytes:    maxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          applogger.Error(),
		BaseContext:       nil,
		ConnContext:       nil,
	}

	return &AppServer{server: server, handlers: appHandler, conf: appConf, AppLogger: applogger}
}

func Run(appServer *AppServer) error {
	go catchSignal(appServer)

	return appServer.server.ListenAndServe()
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
			time.Sleep(timeOutClosedServer)

			appServer.handlers.Pool.Stop()

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
