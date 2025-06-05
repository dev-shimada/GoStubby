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
			name: "url match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URL: "/path",
					},
				},
				gotRawPath: "/path/",
			},
			want: true,
		},
		{
			name: "url does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URL: "/path",
					},
				},
				gotRawPath: "/path/?query=param",
			},
			want: false,
		},
		{
			name: "url pattern match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPattern: "^/[a-zA-Z0-9]{3}$",
					},
				},
				gotRawPath: "/aA0/",
			},
			want: true,
		},
		{
			name: "url pattern does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPattern: "^/[a-zA-Z0-9]{3}$",
					},
				},
				gotRawPath: "/aA01/",
			},
			want: false,
		},
		{
			name: "url path match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPath: "/path/",
					},
				},
				gotPath: "/path/",
			},
			want: true,
		},
		{
			name: "url path does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPath: "/path/",
					},
				},
				gotPath: "/path/path2/",
			},
			want: false,
		},
		{
			name: "url path pattern match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathPattern: "^/[a-zA-Z0-9]{3}$",
					},
				},
				gotPath: "/aA0/",
			},
			want: true,
		},
		{
			name: "url path pattern does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathPattern: "^/[a-zA-Z0-9]{3}$",
					},
				},
				gotPath: "/aA0a/",
			},
			want: false,
		},
		{
			name: "url path template equalTo match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								EqualTo: "12345",
							},
							"param2": {
								EqualTo: "abcde",
							},
						},
					},
				},
				gotPath: "/path/12345/abcde",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "12345",
				"param2": "abcde",
			},
		},
		{
			name: "url path template equalTo does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								EqualTo: "12345",
							},
							"param2": {
								EqualTo: "abcde",
							},
						},
					},
				},
				gotPath: "/path/12345/abcdef",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template matches match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^[0-9]{5}$",
							},
						},
					},
				},
				gotPath: "/path/12345",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "12345",
			},
		},
		{
			name: "url path template matches does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^[0-9]{5}$",
							},
						},
					},
				},
				gotPath: "/path/123456",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^[0-9]{5}$",
							},
						},
					},
				},
				gotPath: "/path/123456",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "123456",
			},
		},
		{
			name: "url path template doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^[0-9]{5}$",
							},
						},
					},
				},
				gotPath: "/path/12345",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Contains: "abc",
							},
						},
					},
				},
				gotPath: "/path/xabcyz",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "xabcyz",
			},
		},
		{
			name: "url path template contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Contains: "abc",
							},
						},
					},
				},
				gotPath: "/path/xyz",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "abc",
							},
						},
					},
				},
				gotPath: "/path/xyz",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "xyz",
			},
		},
		{
			name: "url path template doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "abc",
							},
						},
					},
				},
				gotPath: "/path/xabcyz",
			},
			want:    false,
			wantMap: nil,
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
			name: "empty",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{},
				},
				gotRawQuery: url.Values{},
				gotQuery:    url.Values{},
			},
			want:    true,
			wantMap: map[string]string{},
		},
		{
			name: "empty config",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{},
				},
				gotRawQuery: url.Values{
					"param1": []string{"12345"},
				},
				gotQuery: url.Values{
					"param1": []string{"encoded"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "encoded",
			},
		},
		{
			name: "query equalTo match",
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
		{
			name: "query equalTo does not match",
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
					"param1": []string{"123456"},
				},
				gotQuery: url.Values{
					"param1": []string{"123456"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "query matches match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^[0-9]{5}$",
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
		{
			name: "query matches does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^[0-9]{5}$",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"123456"},
				},
				gotQuery: url.Values{
					"param1": []string{"123456"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "query doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^[0-9]{5}$",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"123456"},
				},
				gotQuery: url.Values{
					"param1": []string{"123456"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "123456",
			},
		},
		{
			name: "query doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^[0-9]{5}$",
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
			want:    false,
			wantMap: nil,
		},
		{
			name: "query contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Contains: "abc",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"xabcyz"},
				},
				gotQuery: url.Values{
					"param1": []string{"xabcyz"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "xabcyz",
			},
		},
		{
			name: "query contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Contains: "abc",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"xyz"},
				},
				gotQuery: url.Values{
					"param1": []string{"xyz"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "query doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "abc",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"xyz"},
				},
				gotQuery: url.Values{
					"param1": []string{"xyz"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "xyz",
			},
		},
		{
			name: "query doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "abc",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"xabcyz"},
				},
				gotQuery: url.Values{
					"param1": []string{"xabcyz"},
				},
			},
			want:    false,
			wantMap: nil,
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

func Test_HeaderMatcher(t *testing.T) {
	type args struct {
		endpoint model.Endpoint
		headers  map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMap map[string][]string
	}{
		{
			name: "empty headers",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{},
				},
				headers: map[string][]string{},
			},
			want:    true,
			wantMap: map[string][]string{},
		},
		{
			name: "header equalTo match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Content-Type": {
								EqualTo: "application/json",
							},
						},
					},
				},
				headers: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			want: true,
			wantMap: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "header equalTo does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Content-Type": {
								EqualTo: "application/json",
							},
						},
					},
				},
				headers: map[string][]string{
					"Content-Type": {"text/plain"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "header matches pattern match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Authorization": {
								Matches: "^Bearer [A-Za-z0-9-_]+\\.[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]*$",
							},
						},
					},
				},
				headers: map[string][]string{
					"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
				},
			},
			want: true,
			wantMap: map[string][]string{
				"Authorization": {
					"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
				},
			},
		},
		{
			name: "header matches pattern does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Authorization": {
								Matches: "^Bearer [A-Za-z0-9-_]+\\.[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]*$",
							},
						},
					},
				},
				headers: map[string][]string{
					"Authorization": {"Basic dXNlcjpwYXNzd29yZA=="},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "header doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Authorization": {
								DoesNotMatch: "^Bearer [A-Za-z0-9-_]+\\.[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]*$",
							},
						},
					},
				},
				headers: map[string][]string{
					"Authorization": {"Basic dXNlcjpwYXNzd29yZA=="},
				},
			},
			want: true,
			wantMap: map[string][]string{
				"Authorization": {"Basic dXNlcjpwYXNzd29yZA=="},
			},
		},
		{
			name: "header doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Authorization": {
								DoesNotMatch: "^Bearer [A-Za-z0-9-_]+\\.[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]*$",
							},
						},
					},
				},
				headers: map[string][]string{
					"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "header contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Accept": {
								Contains: "application/json",
							},
						},
					},
				},
				headers: map[string][]string{
					"Accept": {"text/html, application/json, */*"},
				},
			},
			want: true,
			wantMap: map[string][]string{
				"Accept": {"text/html, application/json, */*"},
			},
		},
		{
			name: "header contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Accept": {
								Contains: "application/json",
							},
						},
					},
				},
				headers: map[string][]string{
					"Accept": {"text/html, text/plain, */*"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "header doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Accept": {
								DoesNotContain: "application/json",
							},
						},
					},
				},
				headers: map[string][]string{
					"Accept": {"text/html, text/plain, */*"},
				},
			},
			want: true,
			wantMap: map[string][]string{
				"Accept": {"text/html, text/plain, */*"},
			},
		},
		{
			name: "header doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Accept": {
								DoesNotContain: "application/json",
							},
						},
					},
				},
				headers: map[string][]string{
					"Accept": {"text/html, application/json, */*"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "multiple headers match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Content-Type": {
								EqualTo: "application/json",
							},
							"Authorization": {
								Contains: "Bearer",
							},
						},
					},
				},
				headers: map[string][]string{
					"Content-Type":  {"application/json"},
					"Authorization": {"Bearer token123"},
				},
			},
			want: true,
			wantMap: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer token123"},
			},
		},
		{
			name: "multiple headers one mismatch",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Headers: map[string]model.Matcher{
							"Content-Type": {
								EqualTo: "application/json",
							},
							"Authorization": {
								Contains: "Bearer",
							},
						},
					},
				},
				headers: map[string][]string{
					"Content-Type":  {"application/xml"},
					"Authorization": {"Bearer token123"},
				},
			},
			want:    false,
			wantMap: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMap := tt.args.endpoint.HeaderMatcher(tt.args.headers); got != tt.want {
				t.Errorf("headerMatcher() = %v, want %v", got, tt.want)
			} else if !cmp.Equal(gotMap, tt.wantMap) {
				t.Errorf("diff: %v", cmp.Diff(gotMap, tt.wantMap))
			}
		})
	}
}

func Test_BodyMatcher(t *testing.T) {
	type args struct {
		endpoint model.Endpoint
		body     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "body equalTo match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							EqualTo: "test",
						},
					},
				},
				body: "test",
			},
			want: true,
		},
		{
			name: "body equalTo does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							EqualTo: "test",
						},
					},
				},
				body: "different",
			},
			want: false,
		},
		{
			name: "body matches match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Matches: "^test[0-9]+$",
						},
					},
				},
				body: "test123",
			},
			want: true,
		},
		{
			name: "body matches does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Matches: "^test[0-9]+$",
						},
					},
				},
				body: "testxyz",
			},
			want: false,
		},
		{
			name: "body doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotMatch: "^test[0-9]+$",
						},
					},
				},
				body: "testxyz",
			},
			want: true,
		},
		{
			name: "body doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotMatch: "^test[0-9]+$",
						},
					},
				},
				body: "test123",
			},
			want: false,
		},
		{
			name: "body contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Contains: "test",
						},
					},
				},
				body: "this is a test message",
			},
			want: true,
		},
		{
			name: "body contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Contains: "test",
						},
					},
				},
				body: "this is a sample",
			},
			want: false,
		},
		{
			name: "body doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotContain: "test",
						},
					},
				},
				body: "this is a sample",
			},
			want: true,
		},
		{
			name: "body doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotContain: "test",
						},
					},
				},
				body: "this is a test message",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.endpoint.BodyMatcher(tt.args.body); got != tt.want {
				t.Errorf("bodyMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}
