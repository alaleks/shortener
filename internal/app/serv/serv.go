// Package serv creates, starts and stops a web server.
package serv

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/alaleks/shortener/proto"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/handlers"
	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/router"
	"github.com/alaleks/shortener/internal/app/serv/middleware"
	"github.com/alaleks/shortener/internal/app/serv/middleware/auth"
	"github.com/alaleks/shortener/internal/app/serv/middleware/compress"
	"github.com/alaleks/shortener/internal/app/storage"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
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
	grpc     *grpc.Server
	handlers *handlers.Handlers
	Logger   *logger.AppLogger
	cfg      config.Configurator
}

// New creates a new server.
func New() *AppServer {
	var (
		cfg        config.Configurator = config.New(config.Options{Env: true, Flag: true})
		logger                         = logger.NewLogger()
		st                             = storage.InitStore(cfg, logger)
		appHandler                     = handlers.New(cfg, logger, st)
		auth                           = auth.TurnOn(appHandler.Storage.St, cfg.GetSecretKey())
		routers                        = router.Create(appHandler)
	)

	server := &http.Server{
		Handler: middleware.New(compress.Compression, compress.Decompression, auth.Authorization).
			Configure(routers),
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              cfg.GetServAddr(),
		ErrorLog:          log.New(logger, "", 0),
		TLSConfig:         nil,
		MaxHeaderBytes:    maxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	// register grpc server.
	grpc := grpc.NewServer()
	pbSrv := pb.New(st, logger, cfg.GetSecretKey(), cfg.GetTrustedSubnet())
	pb.RegisterShortenerServer(grpc, pbSrv)

	return &AppServer{
		server:   server,
		handlers: appHandler,
		cfg:      cfg,
		grpc:     grpc,
		Logger:   logger,
	}
}

// Run starts the server.
func Run(appServer *AppServer) error {
	go catchSignal(appServer)

	// run grpc server
	go func() {
		listener, err := net.Listen("tcp", appServer.cfg.GetGRPCPort())
		if err != nil {
			appServer.Logger.LZ.Error(err)

			return
		}

		if err := appServer.grpc.Serve(listener); err != nil {
			appServer.Logger.LZ.Error(err)

			return
		}
	}()

	// turn on tls for https connections
	if appServer.cfg.EnableTLS() {
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

	_ = http2.ConfigureServer(a.server, &http2.Server{})
}

// catchSignal will catch SIGINT, SIGHUP, SIGQUIT and SIGTERM and close the server.
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
			appServer.grpc.GracefulStop()

			if err := appServer.handlers.Storage.St.Close(); err != nil {
				appServer.Logger.LZ.Fatal(err)
			}

			if err := appServer.server.Shutdown(context.Background()); err != nil {
				appServer.Logger.LZ.Fatal(err)
			}
		case <-reloadSignals:
			if err := appServer.handlers.Storage.St.Close(); err != nil {
				appServer.Logger.LZ.Fatal(err)
			}
		}
	}
}
