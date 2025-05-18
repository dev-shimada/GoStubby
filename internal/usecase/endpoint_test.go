package usecase_test

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/dev-shimada/gostubby/internal/domain/repository"
	"github.com/dev-shimada/gostubby/internal/usecase"
	"github.com/google/go-cmp/cmp"
)

// モックリポジトリ実装
type mockConfigRepository struct {
	endpoints []model.Endpoint
	err       error
}

func (m *mockConfigRepository) Load(path string) ([]model.Endpoint, error) {
	return m.endpoints, m.err
}

func TestEndpointUsecase_EndpointMatcher(t *testing.T) {
	type fields struct {
		cr repository.ConfigRepository
	}
	type args struct {
		arg usecase.EndpointMatcherArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.EndpointMatcherResult
		wantErr bool
	}{
		{
			name: "パスパラメータの完全一致",
			fields: fields{
				cr: &mockConfigRepository{
					endpoints: []model.Endpoint{
						{
							Name: "Test Endpoint",
							Request: model.Request{
								Method:          "GET",
								URLPathTemplate: "/users/{id}",
								PathParameters: map[string]model.Matcher{
									"id": {
										EqualTo: "123",
									},
								},
							},
							Response: model.Response{
								Status: 200,
								Body:   `{"id": "123", "name": "Test User"}`,
							},
						},
					},
				},
			},
			args: args{
				arg: usecase.EndpointMatcherArgs{
					Request: struct {
						UrlRawPath     string
						UrlPath        string
						Body           io.ReadCloser
						Method         string
						RawQueryValues url.Values
						QueryValues    url.Values
					}{
						UrlRawPath: "/users/123",
						UrlPath:    "/users/123",
						Body:       io.NopCloser(strings.NewReader("")),
						Method:     "GET",
					},
					ConfigPath: "test-config.json",
				},
			},
			want: usecase.EndpointMatcherResult{
				Endpoint: model.Endpoint{
					Name: "Test Endpoint",
					Request: model.Request{
						Method:          "GET",
						URLPathTemplate: "/users/{id}",
						PathParameters: map[string]model.Matcher{
							"id": {
								EqualTo: "123",
							},
						},
					},
					Response: model.Response{
						Status: 200,
						Body:   `{"id": "123", "name": "Test User"}`,
					},
				},
				ResponseBody:   `{"id": "123", "name": "Test User"}`,
				ResponseStatus: 200,
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path: map[string]string{
						"id": "123",
					},
					Query: map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name: "正規表現によるパスパラメータのマッチング",
			fields: fields{
				cr: &mockConfigRepository{
					endpoints: []model.Endpoint{
						{
							Name: "Regex Test Endpoint",
							Request: model.Request{
								Method:          "GET",
								URLPathTemplate: "/products/{id}",
								PathParameters: map[string]model.Matcher{
									"id": {
										Matches: "^\\d+$",
									},
								},
							},
							Response: model.Response{
								Status: 200,
								Body:   `{"id": "456", "name": "Test Product"}`,
							},
						},
					},
				},
			},
			args: args{
				arg: usecase.EndpointMatcherArgs{
					Request: struct {
						UrlRawPath     string
						UrlPath        string
						Body           io.ReadCloser
						Method         string
						RawQueryValues url.Values
						QueryValues    url.Values
					}{
						UrlRawPath: "/products/456",
						UrlPath:    "/products/456",
						Body:       io.NopCloser(strings.NewReader("")),
						Method:     "GET",
					},
					ConfigPath: "test-config.json",
				},
			},
			want: usecase.EndpointMatcherResult{
				Endpoint: model.Endpoint{
					Name: "Regex Test Endpoint",
					Request: model.Request{
						Method:          "GET",
						URLPathTemplate: "/products/{id}",
						PathParameters: map[string]model.Matcher{
							"id": {
								Matches: "^\\d+$",
							},
						},
					},
					Response: model.Response{
						Status: 200,
						Body:   `{"id": "456", "name": "Test Product"}`,
					},
				},
				ResponseBody:   `{"id": "456", "name": "Test Product"}`,
				ResponseStatus: 200,
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path: map[string]string{
						"id": "456",
					},
					Query: map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name: "クエリパラメータのマッチング",
			fields: fields{
				cr: &mockConfigRepository{
					endpoints: []model.Endpoint{
						{
							Name: "Query Test Endpoint",
							Request: model.Request{
								Method:          "GET",
								URLPathTemplate: "/search",
								QueryParameters: map[string]model.Matcher{
									"q": {
										Contains: "test",
									},
									"page": {
										Matches: "^\\d+$",
									},
								},
							},
							Response: model.Response{
								Status: 200,
								Body:   `{"results": [{"name": "Test Result"}]}`,
							},
						},
					},
				},
			},
			args: args{
				arg: usecase.EndpointMatcherArgs{
					Request: struct {
						UrlRawPath     string
						UrlPath        string
						Body           io.ReadCloser
						Method         string
						RawQueryValues url.Values
						QueryValues    url.Values
					}{
						UrlRawPath: "/search",
						UrlPath:    "/search",
						Body:       io.NopCloser(strings.NewReader("")),
						Method:     "GET",
						RawQueryValues: url.Values{
							"q":    []string{"test-query"},
							"page": []string{"1"},
						},
						QueryValues: url.Values{
							"q":    []string{"test-query"},
							"page": []string{"1"},
						},
					},
					ConfigPath: "test-config.json",
				},
			},
			want: usecase.EndpointMatcherResult{
				Endpoint: model.Endpoint{
					Name: "Query Test Endpoint",
					Request: model.Request{
						Method:          "GET",
						URLPathTemplate: "/search",
						QueryParameters: map[string]model.Matcher{
							"q": {
								Contains: "test",
							},
							"page": {
								Matches: "^\\d+$",
							},
						},
					},
					Response: model.Response{
						Status: 200,
						Body:   `{"results": [{"name": "Test Result"}]}`,
					},
				},
				ResponseBody:   `{"results": [{"name": "Test Result"}]}`,
				ResponseStatus: 200,
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path: map[string]string{},
					Query: map[string]string{
						"q":    "test-query",
						"page": "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "リクエストボディのマッチング",
			fields: fields{
				cr: &mockConfigRepository{
					endpoints: []model.Endpoint{
						{
							Name: "Body Test Endpoint",
							Request: model.Request{
								Method:          "POST",
								URLPathTemplate: "/api/users",
								Body: model.Matcher{
									Contains:       "email",
									DoesNotContain: "password",
								},
							},
							Response: model.Response{
								Status: 201,
								Body:   `{"id": "789", "status": "created"}`,
							},
						},
					},
				},
			},
			args: args{
				arg: usecase.EndpointMatcherArgs{
					Request: struct {
						UrlRawPath     string
						UrlPath        string
						Body           io.ReadCloser
						Method         string
						RawQueryValues url.Values
						QueryValues    url.Values
					}{
						UrlRawPath: "/api/users",
						UrlPath:    "/api/users",
						Body:       io.NopCloser(strings.NewReader(`{"name":"John","email":"john@example.com"}`)),
						Method:     "POST",
					},
					ConfigPath: "test-config.json",
				},
			},
			want: usecase.EndpointMatcherResult{
				Endpoint: model.Endpoint{
					Name: "Body Test Endpoint",
					Request: model.Request{
						Method:          "POST",
						URLPathTemplate: "/api/users",
						Body: model.Matcher{
							Contains:       "email",
							DoesNotContain: "password",
						},
					},
					Response: model.Response{
						Status: 201,
						Body:   `{"id": "789", "status": "created"}`,
					},
				},
				ResponseBody:   `{"id": "789", "status": "created"}`,
				ResponseStatus: 201,
				Data: struct {
					Path  map[string]string
					Query map[string]string
				}{
					Path:  map[string]string{},
					Query: map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name: "マッチするエンドポイントがない場合",
			fields: fields{
				cr: &mockConfigRepository{
					endpoints: []model.Endpoint{
						{
							Name: "Test Endpoint",
							Request: model.Request{
								Method:          "GET",
								URLPathTemplate: "/users/{id}",
								PathParameters: map[string]model.Matcher{
									"id": {
										EqualTo: "123",
									},
								},
							},
							Response: model.Response{
								Status: 200,
								Body:   `{"id": "123", "name": "Test User"}`,
							},
						},
					},
				},
			},
			args: args{
				arg: usecase.EndpointMatcherArgs{
					Request: struct {
						UrlRawPath     string
						UrlPath        string
						Body           io.ReadCloser
						Method         string
						RawQueryValues url.Values
						QueryValues    url.Values
					}{
						UrlRawPath: "/users/456", // 存在しないID
						UrlPath:    "/users/456",
						Body:       io.NopCloser(strings.NewReader("")),
						Method:     "GET",
					},
					ConfigPath: "test-config.json",
				},
			},
			want:    usecase.EndpointMatcherResult{},
			wantErr: true,
		},
		{
			name: "リポジトリエラー",
			fields: fields{
				cr: &mockConfigRepository{
					endpoints: nil,
					err:       fmt.Errorf("設定ファイルの読み込みエラー"),
				},
			},
			args: args{
				arg: usecase.EndpointMatcherArgs{
					Request: struct {
						UrlRawPath     string
						UrlPath        string
						Body           io.ReadCloser
						Method         string
						RawQueryValues url.Values
						QueryValues    url.Values
					}{
						UrlRawPath: "/users/123",
						UrlPath:    "/users/123",
						Body:       io.NopCloser(strings.NewReader("")),
						Method:     "GET",
					},
					ConfigPath: "test-config.json",
				},
			},
			want:    usecase.EndpointMatcherResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eu := usecase.NewEndpointUsecase(tt.fields.cr)
			got, err := eu.EndpointMatcher(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndpointUsecase.EndpointMatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				diff := cmp.Diff(tt.want, got)
				if diff != "" {
					t.Errorf("EndpointUsecase.EndpointMatcher() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestEndpointUsecase_ResponseCreator(t *testing.T) {
	type fields struct {
		cr repository.ConfigRepository
	}
	type args struct {
		arg usecase.ResponseCreatorArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.ResponseCreatorResult
		wantErr bool
	}{
		{
			name: "有効なテンプレートの解析",
			fields: fields{
				cr: &mockConfigRepository{},
			},
			args: args{
				arg: usecase.ResponseCreatorArgs{
					Request: struct {
						UrlQuery url.Values
					}{
						UrlQuery: url.Values{
							"param": []string{"value"},
						},
					},
					Endpoint: model.Endpoint{
						Name: "Template Test Endpoint",
						Response: model.Response{
							Status: 200,
							Body:   `{"message": "Hello, world!"}`,
						},
					},
					ResponseBody: `{"id": "{{.Path.id}}", "query": "{{.Query.param}}"}`,
					PathMap: map[string]string{
						"id": "123",
					},
				},
			},
			want: func() usecase.ResponseCreatorResult {
				tpl, _ := template.New("response").Parse(`{"id": "{{.Path.id}}", "query": "{{.Query.param}}"}`)
				return usecase.ResponseCreatorResult{
					Template: tpl,
				}
			}(),
			wantErr: false,
		},
		{
			name: "無効なテンプレート構文",
			fields: fields{
				cr: &mockConfigRepository{},
			},
			args: args{
				arg: usecase.ResponseCreatorArgs{
					ResponseBody: `{"id": "{{.Path.id", "invalid": "{{.Invalid}"}`,
					PathMap: map[string]string{
						"id": "123",
					},
				},
			},
			want:    usecase.ResponseCreatorResult{},
			wantErr: true,
		},
		{
			name: "空のテンプレート",
			fields: fields{
				cr: &mockConfigRepository{},
			},
			args: args{
				arg: usecase.ResponseCreatorArgs{
					ResponseBody: "",
				},
			},
			want: func() usecase.ResponseCreatorResult {
				tpl, _ := template.New("response").Parse("")
				return usecase.ResponseCreatorResult{
					Template: tpl,
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eu := usecase.NewEndpointUsecase(tt.fields.cr)
			got, err := eu.ResponseCreator(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndpointUsecase.ResponseCreator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// テンプレートを直接比較するのではなく、実行結果を比較
				// テンプレート変数用のモックデータ
				type templateData struct {
					Path  map[string]string
					Query map[string]string
				}
				mockData := templateData{
					Path: map[string]string{
						"id": "123",
					},
					Query: map[string]string{
						"param": "value",
					},
				}

				// 両方のテンプレートを同じモックデータで実行して結果を比較
				var gotBuf, wantBuf strings.Builder
				if got.Template != nil {
					if err := got.Template.Execute(&gotBuf, mockData); err != nil {
						t.Errorf("Failed to execute got template: %v", err)
						return
					}
				}
				if tt.want.Template != nil {
					if err := tt.want.Template.Execute(&wantBuf, mockData); err != nil {
						t.Errorf("Failed to execute want template: %v", err)
						return
					}
				}
				gotText := gotBuf.String()
				wantText := wantBuf.String()

				if !cmp.Equal(gotText, wantText) {
					t.Errorf("Template strings differ:\n%s", cmp.Diff(wantText, gotText))
				}
			}
		})
	}
}
