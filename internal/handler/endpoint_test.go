package handler_test

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/dev-shimada/gostubby/internal/handler"
	"github.com/dev-shimada/gostubby/internal/usecase"
)

type mockEndpointUsecase struct {
	endpointMatcherFunc func(usecase.EndpointMatcherArgs) (usecase.EndpointMatcherResult, error)
	responseCreatorFunc func(usecase.ResponseCreatorArgs) (usecase.ResponseCreatorResult, error)
}

func (m *mockEndpointUsecase) EndpointMatcher(args usecase.EndpointMatcherArgs) (usecase.EndpointMatcherResult, error) {
	return m.endpointMatcherFunc(args)
}

func (m *mockEndpointUsecase) ResponseCreator(args usecase.ResponseCreatorArgs) (usecase.ResponseCreatorResult, error) {
	return m.responseCreatorFunc(args)
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		request        *http.Request
		matcherResult  usecase.EndpointMatcherResult
		matcherErr     error
		creatorResult  usecase.ResponseCreatorResult
		creatorErr     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Successful request",
			configPath: "test/config.json",
			request:    httptest.NewRequest(http.MethodGet, "/test?param=value", nil),
			matcherResult: usecase.EndpointMatcherResult{
				Endpoint: model.Endpoint{
					Name: "test-endpoint",
				},
				ResponseStatus: http.StatusOK,
				ResponseBody:   "template content",
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path:  map[string]string{"id": "123"},
					Query: map[string]string{"param": "value"},
				},
			},
			matcherErr: nil,
			creatorResult: usecase.ResponseCreatorResult{
				Template: template.Must(template.New("test").Parse("Hello {{.Path.id}}")),
			},
			creatorErr:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello 123",
		},
		{
			name:           "EndpointMatcher error",
			configPath:     "test/config.json",
			request:        httptest.NewRequest(http.MethodGet, "/test", nil),
			matcherErr:     &url.Error{Op: "parse", URL: "invalid", Err: url.InvalidHostError("invalid")},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "ResponseCreator error",
			configPath: "test/config.json",
			request:    httptest.NewRequest(http.MethodGet, "/test", nil),
			matcherResult: usecase.EndpointMatcherResult{
				Endpoint: model.Endpoint{
					Name: "test-endpoint",
				},
				ResponseStatus: http.StatusNotFound,
				ResponseBody:   "template content",
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path:  map[string]string{},
					Query: map[string]string{},
				},
			},
			matcherErr:     nil,
			creatorResult:  usecase.ResponseCreatorResult{},
			creatorErr:     errors.New("template parsing error"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &mockEndpointUsecase{
				endpointMatcherFunc: func(args usecase.EndpointMatcherArgs) (usecase.EndpointMatcherResult, error) {
					if tt.matcherErr != nil {
						return usecase.EndpointMatcherResult{}, tt.matcherErr
					}
					return tt.matcherResult, nil
				},
				responseCreatorFunc: func(args usecase.ResponseCreatorArgs) (usecase.ResponseCreatorResult, error) {
					if tt.creatorErr != nil {
						return usecase.ResponseCreatorResult{}, tt.creatorErr
					}
					return tt.creatorResult, nil
				},
			}

			handler := handler.NewEndpointHandler(tt.configPath, mockUsecase)
			w := httptest.NewRecorder()
			handler.Handle(w, tt.request)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				if body := w.Body.String(); body != tt.expectedBody {
					t.Errorf("Expected body %q, got %q", tt.expectedBody, body)
				}
			}
		})
	}
}

func Test_rawQueryValues(t *testing.T) {
	type args struct {
		r http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "正常なクエリパラメータ",
			args: args{
				r: http.Request{
					URL: &url.URL{
						RawQuery: "param1=value1&param2=value2",
					},
				},
			},
			want: url.Values{
				"param1": []string{"value1"},
				"param2": []string{"value2"},
			},
			wantErr: false,
		},
		{
			name: "空のクエリパラメータ",
			args: args{
				r: http.Request{
					URL: &url.URL{
						RawQuery: "",
					},
				},
			},
			want:    url.Values{},
			wantErr: false,
		},
		{
			name: "不正な形式のクエリパラメータ",
			args: args{
				r: http.Request{
					URL: &url.URL{
						RawQuery: "param1=value1&param2",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 元の実装ではなく、テスト用の実装を使用
			got, err := handler.ExportedRawQueryValues(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("rawQueryValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rawQueryValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
