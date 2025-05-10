package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/dev-shimada/gostubby/internal/domain/repository"
	"github.com/dev-shimada/gostubby/internal/infrastructure/config"
)

var (
	configPath string
)

func main() {
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

	mux.HandleFunc("/", handle)

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

type endpointUsecase interface {
	endpointMatcher(EndpointMatcherArgs) (EndpointMatcherResult, error)
	responseCreator(ResponseCreatorArgs) (ResponseCreatorResult, error)
}

func handle(w http.ResponseWriter, r *http.Request) {
	cr := config.NewConfigRepository()
	var ne endpointUsecase = NewEndpointUsecase(cr)
	rqv, err := rawQueryValues(*r)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse query parameters: %s", err))
		http.NotFound(w, r)
		return
	}

	EndpointMatcherArgs := EndpointMatcherArgs{
		Request: struct {
			UrlRawPath     string
			UrlPath        string
			Body           io.ReadCloser
			Method         string
			RawQueryValues url.Values
			QueryValues    url.Values
		}{
			UrlRawPath:     r.URL.RawPath,
			UrlPath:        r.URL.Path,
			Body:           r.Body,
			Method:         r.Method,
			RawQueryValues: rqv,
			QueryValues:    r.URL.Query(),
		},
		ConfigPath: configPath,
	}
	em, err := ne.endpointMatcher(EndpointMatcherArgs)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to match endpoint: %v", err))
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(em.ResponseStatus)
	ResponseCreatorArgs := ResponseCreatorArgs{
		Request: struct {
			UrlQuery url.Values
		}{
			UrlQuery: r.URL.Query(),
		},
		Endpoint:     em.Endpoint,
		ResponseBody: em.ResponseBody,
	}
	rc, err := ne.responseCreator(ResponseCreatorArgs)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create response: %v", err))
		http.NotFound(w, r)
		return
	}
	if err := rc.Template.Execute(w, em.Data); err != nil {
		slog.Error(fmt.Sprintf("Failed to execute template: %s", err))
		http.NotFound(w, r)
		return
	}
}

// rawQueryValues parses the raw query string from the request URL and returns a url.Values map.
// It splits the query string by '&' and then splits each key-value pair by '='.
// If the query string is malformed, it returns an error.
func rawQueryValues(r http.Request) (url.Values, error) {
	ret := url.Values{}
	for v := range strings.SplitSeq(r.URL.RawQuery, "&") {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid query parameter: %s", v)
		}
		ret.Add(kv[0], kv[1])
	}
	return ret, nil
}

type EndpointUsecase struct {
	cr repository.ConfigRepository
}

func NewEndpointUsecase(cr repository.ConfigRepository) *EndpointUsecase {
	return &EndpointUsecase{
		cr: cr,
	}
}

type EndpointMatcherArgs struct {
	Request struct {
		UrlRawPath     string
		UrlPath        string
		Body           io.ReadCloser
		Method         string
		RawQueryValues url.Values
		QueryValues    url.Values
	}
	ConfigPath string
}
type EndpointMatcherResult struct {
	Endpoint       model.Endpoint
	ResponseBody   string
	ResponseStatus int
	Data           struct {
		Path  map[string]string
		Query map[string]string
	}
}

func (eu EndpointUsecase) endpointMatcher(arg EndpointMatcherArgs) (EndpointMatcherResult, error) {
	endpoints, err := eu.cr.Load(configPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		return EndpointMatcherResult{}, err
	}
	for _, e := range endpoints {
		var responseBody string
		switch {
		case e.Response.BodyFileName != "":
			file, err := os.Open(e.Response.BodyFileName)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to open body file: %s", err))
				return EndpointMatcherResult{}, err
			}
			defer func() {
				if err := file.Close(); err != nil {
					slog.Error(fmt.Sprintf("Failed to close file: %s", err))
				}
			}()
			body, err := io.ReadAll(file)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to read body file: %s", err))
				return EndpointMatcherResult{}, err
			}
			responseBody = string(body)
		case e.Response.Body != "":
			responseBody = e.Response.Body
		default:
			slog.Error("Response body is empty")
			return EndpointMatcherResult{}, fmt.Errorf("response body is empty")
		}

		// pathMatcher
		isMatchPath, pathMap := e.PathMatcher(arg.Request.UrlRawPath, arg.Request.UrlPath)
		// queryMatcher
		isMatchQuery, queryMap := e.QueryMatcher(arg.Request.RawQueryValues, arg.Request.QueryValues)
		body, err := io.ReadAll(arg.Request.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read request body: %s", err))
			return EndpointMatcherResult{}, err
		}
		isMatchBody := e.BodyMatcher(string(body))
		if arg.Request.Method == e.Request.Method && isMatchPath && isMatchQuery && isMatchBody {
			slog.Info(fmt.Sprintf("Matched endpoint: %s", e.Name))
			return EndpointMatcherResult{
				Endpoint:       e,
				ResponseBody:   responseBody,
				ResponseStatus: e.Response.Status,
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path:  pathMap,
					Query: queryMap,
				},
			}, nil
		}
	}
	return EndpointMatcherResult{}, fmt.Errorf("no matching endpoint found")
}

type ResponseCreatorArgs struct {
	Request struct {
		UrlQuery url.Values
	}
	Endpoint     model.Endpoint
	ResponseBody string
	PathMap      map[string]string
}
type ResponseCreatorResult struct {
	Template *template.Template
}

func (eu EndpointUsecase) responseCreator(arg ResponseCreatorArgs) (ResponseCreatorResult, error) {
	tpl, err := template.New("response").Parse(arg.ResponseBody)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse response template: %s", err))
		return ResponseCreatorResult{}, err
	}
	return ResponseCreatorResult{
		Template: tpl,
	}, nil
}
