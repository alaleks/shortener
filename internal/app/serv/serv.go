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
	Logger   *logger.AppLogger
	conf     config.Configurator
}

func New(sizeUID int) *AppServer {
	var (
		appConf    config.Configurator = config.New(config.Options{Env: true, Flag: true}, sizeUID)
		logger                         = logger.NewLogger()
		appHandler                     = handlers.New(appConf, logger)
		auth                           = auth.TurnOn(appHandler.Storage, appConf.GetSecretKey())
		routers                        = router.Create(appHandler)
	)

	server := &http.Server{
		Handler: middleware.New(compress.Compression, compress.Unpacking, auth.Authorization).
			Configure(routers),
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              appConf.GetServAddr(),
		ErrorLog:          log.New(logger, "", 0),
		TLSConfig:         nil,
		MaxHeaderBytes:    maxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	return &AppServer{
		server:   server,
		handlers: appHandler,
		conf:     appConf,
		Logger:   logger,
	}
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
			appServer.handlers.Storage.Pool.Stop()

			if err := appServer.handlers.Storage.Store.Close(); err != nil {
				appServer.Logger.LZ.Fatal(err)
			}

			if err := appServer.server.Shutdown(context.Background()); err != nil {
				appServer.Logger.LZ.Fatal(err)
			}
		case <-reloadSignals:
			if err := appServer.handlers.Storage.Store.Close(); err != nil {
				appServer.Logger.LZ.Fatal(err)
			}
		}
	}
}
