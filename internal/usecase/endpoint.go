package usecase

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/url"
	"os"

	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/dev-shimada/gostubby/internal/domain/repository"
)

type EndpointUsecase struct {
	cr repository.ConfigRepository
}

func NewEndpointUsecase(cr repository.ConfigRepository) EndpointUsecase {
	return EndpointUsecase{
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

func (eu EndpointUsecase) EndpointMatcher(arg EndpointMatcherArgs) (EndpointMatcherResult, error) {
	endpoints, err := eu.cr.Load(arg.ConfigPath)
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

func (eu EndpointUsecase) ResponseCreator(arg ResponseCreatorArgs) (ResponseCreatorResult, error) {
	tpl, err := template.New("response").Parse(arg.ResponseBody)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse response template: %s", err))
		return ResponseCreatorResult{}, err
	}
	return ResponseCreatorResult{
		Template: tpl,
	}, nil
}
