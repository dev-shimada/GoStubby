package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-shimada/gostubby/internal/handler"
	"github.com/dev-shimada/gostubby/internal/infrastructure/config"
	"github.com/dev-shimada/gostubby/internal/usecase"
)

var (
	configPath string
)

func main() {
	// json logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))

	var (
		host      string
		port      int
		httpsPort int
		certFile  string
		keyFile   string
		// configPath string
	)
	// Host configuration
	host = *flag.String("h", "localhost", "Host address to bind to (use 0.0.0.0 for Docker)")
	flag.StringVar(&host, "host", "localhost", "Host address to bind to (use 0.0.0.0 for Docker)")

	// HTTP configuration
	port = *flag.Int("p", 8080, "HTTP port number to listen on")
	flag.IntVar(&port, "port", 8080, "HTTP port number to listen on")

	// HTTPS configuration
	httpsPort = *flag.Int("s", 8443, "HTTPS port number to listen on")
	flag.IntVar(&httpsPort, "https-port", 8443, "HTTPS port number to listen on")
	certFile = *flag.String("t", "", "Path to SSL/TLS certificate file")
	flag.StringVar(&certFile, "cert", "", "Path to SSL/TLS certificate file")
	keyFile = *flag.String("k", "", "Path to SSL/TLS private key file")
	flag.StringVar(&keyFile, "key", "", "Path to SSL/TLS private key file")

	// General configuration
	configPath = *flag.String("config", "configs", "Path to configuration directory or file")
	flag.StringVar(&configPath, "c", "configs", "Path to configuration directory or file")
	flag.Parse()

	mux := http.NewServeMux()

	// Dependency injection
	cr := config.NewConfigRepository()
	eu := usecase.NewEndpointUsecase(cr)
	eh := handler.NewEndpointHandler(configPath, *eu)

	mux.HandleFunc("/", eh.Handle)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

	// Create HTTP server
	httpAddr := fmt.Sprintf("%s:%d", host, port)
	httpSrv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	// Create HTTPS server if certificate and key files are provided
	var httpsSrv *http.Server
	if certFile != "" && keyFile != "" {
		httpsAddr := fmt.Sprintf("%s:%d", host, httpsPort)
		httpsSrv = &http.Server{
			Addr:    httpsAddr,
			Handler: mux,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
	}

	// Start HTTP server
	slog.Info(fmt.Sprintf("HTTP server is running at http://%s:%d", host, port))
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("HTTP server closed")
			} else {
				slog.Error(fmt.Sprintf("HTTP ListenAndServe: %v", err))
			}
		}
	}()

	// Start HTTPS server if configured
	if httpsSrv != nil {
		slog.Info(fmt.Sprintf("HTTPS server is running at https://%s:%d", host, httpsPort))
		go func() {
			if err := httpsSrv.ListenAndServeTLS(certFile, keyFile); err != nil {
				if err == http.ErrServerClosed {
					slog.Info("HTTPS server closed")
				} else {
					slog.Error(fmt.Sprintf("HTTPS ListenAndServeTLS: %v", err))
				}
			}
		}()
	}

	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("HTTP server Shutdown: %v", err))
	}

	// Shutdown HTTPS server if it was started
	if httpsSrv != nil {
		if err := httpsSrv.Shutdown(shutdownCtx); err != nil {
			slog.Error(fmt.Sprintf("HTTPS server Shutdown: %v", err))
		}
	}
}
