package model_test

import (
	"net/url"
	"testing"

	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/google/go-cmp/cmp"
)

func Test_PathMatcher(t *testing.T) {
	type args struct {
		endpoint   model.Endpoint
		gotRawPath string
		gotPath    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMap map[string]string
	}{
		{
			name: "url ok",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URL: "http://example.com/path",
					},
				},
				gotRawPath: "http://example.com/path/",
			},
			want: true,
		},
		{
			name: "url ng",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URL: "http://example.com/path",
					},
				},
				gotRawPath: "http://example.com/path/?query=param",
			},
			want: false,
		},
		{
			name: "url pattern ok",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPattern: "^http://example.com/[a-zA-Z0-9]{3}$",
					},
				},
				gotRawPath: "http://example.com/aA0/",
			},
			want: true,
		},
		{
			name: "url pattern ng",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPattern: "^http://example.com/[a-zA-Z0-9]{3}$",
					},
				},
				gotRawPath: "http://example.com/aA01/",
			},
			want: false,
		},
		{
			name: "url path ok",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPath: "http://example.com/path/",
					},
				},
				gotPath: "http://example.com/path/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMap := tt.args.endpoint.PathMatcher(tt.args.gotRawPath, tt.args.gotPath); got != tt.want {
				t.Errorf("pathMatcher() = %v, want %v", got, tt.want)
			} else if !cmp.Equal(gotMap, tt.wantMap) {
				t.Errorf("diff: %v", cmp.Diff(gotMap, tt.wantMap))
			}
		})
	}
}

func Test_QueryMatcher(t *testing.T) {
	type args struct {
		endpoint    model.Endpoint
		gotRawQuery url.Values
		gotQuery    url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMap map[string]string
	}{
		{
			name: "all",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								EqualTo: "12345",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"12345"},
				},
				gotQuery: url.Values{
					"param1": []string{"12345"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMap := tt.args.endpoint.QueryMatcher(tt.args.gotRawQuery, tt.args.gotQuery); got != tt.want {
				t.Errorf("queryMatcher() = %v, want %v", got, tt.want)
			} else if !cmp.Equal(gotMap, tt.wantMap) {
				t.Errorf("diff: %v", cmp.Diff(gotMap, tt.wantMap))
			}
		})
	}
}

func Test_BodyMatcher(t *testing.T) {
	type args struct {
		endpoint model.Endpoint
		gotBody  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							EqualTo: "12345",
						},
					},
				},
				gotBody: "12345",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.endpoint.BodyMatcher(tt.args.gotBody); got != tt.want {
				t.Errorf("bodyMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}
