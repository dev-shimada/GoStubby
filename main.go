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
	"regexp"
	"slices"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/dev-shimada/gostubby/internal/domain/model"
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

// func (c config) handle(w http.ResponseWriter, r *http.Request) {
func handle(w http.ResponseWriter, r *http.Request) {
	config := config.NewConfigRepository()
	endpoints, err := config.Load(configPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	for _, endpoint := range endpoints {
		var responseBody string
		switch {
		case endpoint.Response.BodyFileName != "":
			file, err := os.Open(endpoint.Response.BodyFileName)
			if err != nil {
				http.Error(w, "Failed to open body file", http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := file.Close(); err != nil {
					slog.Error(fmt.Sprintf("Failed to close file: %s", err))
				}
			}()
		case endpoint.Response.Body != "":
			responseBody = endpoint.Response.Body
		default:
			slog.Error("Response body is empty")
			http.Error(w, "Response body is empty", http.StatusInternalServerError)
			return
		}

		// pathMatcher
		isMatchPath, pathMap := pathMatcher(endpoint, r.URL.RawPath, r.URL.Path)
		// queryMatcher
		rqv, err := rawQueryValues(*r)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to parse query parameters: %s", err))
			http.Error(w, "Failed to parse query parameters", http.StatusInternalServerError)
			return
		}
		isMatchQuery := queryMatcher(endpoint, rqv)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read request body: %s", err))
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		isMatchBody := bodyMatcher(endpoint, string(body))
		if r.Method == endpoint.Request.Method && isMatchPath && isMatchQuery && isMatchBody {
			slog.Info(fmt.Sprintf("Matched endpoint: %s", endpoint.Name))
			w.WriteHeader(endpoint.Response.Status)

			type gotParams struct {
				Path  map[string]string
				Query map[string]string
			}
			q := make(map[string]string)
			for k, v := range r.URL.Query() {
				q[k] = v[0]
			}
			gp := gotParams{
				Query: q,
				Path:  pathMap,
			}

			// bodyFileNameが指定されている場合は、bodyは無視される
			if endpoint.Response.BodyFileName != "" {
				file, err := os.Open(endpoint.Response.BodyFileName)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to open body file: %s", err))
					http.Error(w, "Failed to open body file", http.StatusInternalServerError)
					return
				}
				defer func() {
					if err := file.Close(); err != nil {
						slog.Error(fmt.Sprintf("Failed to close file: %s", err))
					}
				}()
				body, err := io.ReadAll(file)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to read body file: %s", err))
					http.Error(w, "Failed to read body file", http.StatusInternalServerError)
					return
				}
				responseBody = string(body)
				tpl, err := template.New("response").Parse(responseBody)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to parse response template: %s", err))
					http.Error(w, "Failed to parse response template", http.StatusInternalServerError)
					return
				}
				if err := tpl.Execute(w, gp); err != nil {
					slog.Error(fmt.Sprintf("Failed to execute response template: %s", err))
					http.Error(w, "Failed to execute response template", http.StatusInternalServerError)
					return
				}
				return
			}

			// bodyFileNameが指定されていない場合は、bodyを使用する
			tpl, err := template.New("response").Parse(responseBody)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to parse response template: %s", err))
				http.Error(w, "Failed to parse response template", http.StatusInternalServerError)
				return
			}
			if err := tpl.Execute(w, gp); err != nil {
				slog.Error(fmt.Sprintf("Failed to execute response template: %s", err))
				http.Error(w, "Failed to execute response template", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.NotFound(w, r)
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

func pathMatcher(endpoint model.Endpoint, gotRawPath, gotPath string) (bool, map[string]string) {
	// trim trailing slashes
	gotPath = strings.TrimRight(gotPath, "/")
	gotRawPath = strings.TrimRight(gotRawPath, "/")

	var url string
	switch {
	case endpoint.Request.URL != "":
		url = strings.TrimRight(endpoint.Request.URL, "/")
		if gotRawPath != url {
			return false, nil
		}
		return true, nil
	case endpoint.Request.URLPattern != "":
		url = strings.TrimRight(endpoint.Request.URLPattern, "/")
		if !regexp.MustCompile(url).MatchString(gotRawPath) {
			return false, nil
		}
		return true, nil
	case endpoint.Request.URLPath != "":
		url = strings.TrimRight(endpoint.Request.URLPath, "/")
		if gotPath != url {
			return false, nil
		}
		return true, nil
	case endpoint.Request.URLPathPattern != "":
		url = strings.TrimRight(endpoint.Request.URLPathPattern, "/")
		if !regexp.MustCompile(url).MatchString(gotPath) {
			return false, nil
		}
		return true, nil
	case endpoint.Request.URLPathTemplate != "":
		url = strings.TrimRight(endpoint.Request.URLPathTemplate, "/")
	default:
		return false, nil
	}

	// check if the path parameters match
	requredPathUnits := strings.Split(url, "/")
	gotPathUnits := strings.Split(gotPath, "/")
	if len(requredPathUnits) != len(gotPathUnits) {
		return false, nil
	}

	// placeholder->position
	posMap := make(map[string]int)
	for k := range endpoint.Request.PathParameters {
		placeHolder := fmt.Sprintf("{%s}", k)
		if i := slices.Index(requredPathUnits, placeHolder); i == -1 {
			slog.Error(fmt.Sprintf("Path parameter %s not found in path %s", k, gotPath))
			return false, nil
		} else {
			posMap[k] = i
		}
	}

	for k, v := range endpoint.Request.PathParameters {
		if v.EqualTo != nil {
			if gotPathUnits[posMap[k]] != fmt.Sprint(v.EqualTo) {
				return false, nil
			}
		}
		if v.Matches != nil {
			if !regexp.MustCompile(v.Matches.(string)).MatchString(gotPathUnits[posMap[k]]) {
				return false, nil
			}
		}
		if v.DoesNotMatch != nil {
			if regexp.MustCompile(v.DoesNotMatch.(string)).MatchString(gotPathUnits[posMap[k]]) {
				return false, nil
			}
		}
		if v.Contains != nil {
			if !strings.Contains(gotPathUnits[posMap[k]], v.Contains.(string)) {
				return false, nil
			}
		}
		if v.DoesNotContain != nil {
			if strings.Contains(gotPathUnits[posMap[k]], v.DoesNotContain.(string)) {
				return false, nil
			}
		}
	}
	ret := make(map[string]string)
	for k, v := range posMap {
		ret[k] = gotPathUnits[v]
	}
	return true, ret
}

func queryMatcher(endpoint model.Endpoint, gotQuery url.Values) bool {
	for k, v := range endpoint.Request.QueryParameters {
		if v.EqualTo != nil {
			if gotQuery.Get(k) != fmt.Sprint(v.EqualTo) {
				return false
			}
		}
		if v.Matches != nil {
			if !regexp.MustCompile(v.Matches.(string)).MatchString(gotQuery.Get(k)) {
				return false
			}
		}
		if v.DoesNotMatch != nil {
			if regexp.MustCompile(v.DoesNotMatch.(string)).MatchString(gotQuery.Get(k)) {
				return false
			}
		}
		if v.Contains != nil {
			if !strings.Contains(gotQuery.Get(k), v.Contains.(string)) {
				return false
			}
		}
		if v.DoesNotContain != nil {
			if strings.Contains(gotQuery.Get(k), v.DoesNotContain.(string)) {
				return false
			}
		}
	}
	return true
}

func bodyMatcher(endpoint model.Endpoint, body string) bool {
	if endpoint.Request.Body.EqualTo != nil {
		if body != fmt.Sprint(endpoint.Request.Body.EqualTo) {
			return false
		}
	}
	if endpoint.Request.Body.Matches != nil {
		if !regexp.MustCompile(endpoint.Request.Body.Matches.(string)).MatchString(body) {
			return false
		}
	}
	if endpoint.Request.Body.DoesNotMatch != nil {
		if regexp.MustCompile(endpoint.Request.Body.DoesNotMatch.(string)).MatchString(body) {
			return false
		}
	}
	if endpoint.Request.Body.Contains != nil {
		if !strings.Contains(body, endpoint.Request.Body.Contains.(string)) {
			return false
		}
	}
	if endpoint.Request.Body.DoesNotContain != nil {
		if strings.Contains(body, endpoint.Request.Body.DoesNotContain.(string)) {
			return false
		}
	}
	return true
}
