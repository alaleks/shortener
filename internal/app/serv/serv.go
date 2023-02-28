// Package serv creates, starts and stops a web server.
package serv

import (
	"context"
	"crypto/tls"
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
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

const (
	defaultTimeout           = 2 * time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultIdleTimeout       = 15 * time.Second
	maxHeaderBytes           = 4096
)

// AppServer represents an application server instance.
type AppServer struct {
	server   *http.Server
	handlers *handlers.Handlers
	Logger   *logger.AppLogger
	conf     config.Configurator
}

// New creates a new server.
func New() *AppServer {
	var (
		appConf    config.Configurator = config.New(config.Options{Env: true, Flag: true})
		logger                         = logger.NewLogger()
		appHandler                     = handlers.New(appConf, logger)
		auth                           = auth.TurnOn(appHandler.Storage, appConf.GetSecretKey())
		routers                        = router.Create(appHandler)
	)

	server := &http.Server{
		Handler: middleware.New(compress.Compression, compress.Decompression, auth.Authorization).
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

// Run starts the server.
func Run(appServer *AppServer) error {
	go catchSignal(appServer)

	// turn on tls for https connections
	if appServer.conf.EnableTLS() {
		appServer.turnOnTLS()
		return appServer.server.ListenAndServeTLS("", "")
	}

	return appServer.server.ListenAndServe()
}

func (a *AppServer) turnOnTLS() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("cert"),
		HostPolicy: autocert.HostWhitelist(a.server.Addr),
	}

	a.server.TLSConfig = &tls.Config{
		GetCertificate:           certManager.GetCertificate,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	a.server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	http2.ConfigureServer(a.server, &http2.Server{})
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
