package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/dev-shimada/gostubby/internal/usecase"
)

type endpointHandler struct {
	configPath string
	eu         endpointUsecase
}

func NewEndpointHandler(configPath string, eu usecase.EndpointUsecase) endpointHandler {
	return endpointHandler{
		configPath: configPath,
		eu:         eu,
	}
}

type endpointUsecase interface {
	EndpointMatcher(usecase.EndpointMatcherArgs) (usecase.EndpointMatcherResult, error)
	ResponseCreator(usecase.ResponseCreatorArgs) (usecase.ResponseCreatorResult, error)
}

func (eh endpointHandler) Handle(w http.ResponseWriter, r *http.Request) {
	configPath := eh.configPath
	rqv, err := rawQueryValues(*r)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse query parameters: %s", err))
		http.NotFound(w, r)
		return
	}

	EndpointMatcherArgs := usecase.EndpointMatcherArgs{
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
	em, err := eh.eu.EndpointMatcher(EndpointMatcherArgs)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to match endpoint: %v", err))
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(em.ResponseStatus)
	ResponseCreatorArgs := usecase.ResponseCreatorArgs{
		Request: struct {
			UrlQuery url.Values
		}{
			UrlQuery: r.URL.Query(),
		},
		Endpoint:     em.Endpoint,
		ResponseBody: em.ResponseBody,
	}
	rc, err := eh.eu.ResponseCreator(ResponseCreatorArgs)
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
	if r.URL.RawQuery == "" {
		return ret, nil
	}
	for v := range strings.SplitSeq(r.URL.RawQuery, "&") {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid query parameter: %s", v)
		}
		ret.Add(kv[0], kv[1])
	}
	return ret, nil
}
