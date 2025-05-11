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
				gotPath: "/path/12345/abcd",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template matches match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^\\d{5}$",
							},
							"param2": {
								Matches: "^[a-z]{5}$",
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
			name: "url path template matches does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^\\d{5}$",
							},
							"param2": {
								Matches: "^[a-z]{5}$",
							},
						},
					},
				},
				gotPath: "/path/12345/abcd1",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^\\d{5}$",
							},
							"param2": {
								DoesNotMatch: "^[a-z]{5}$",
							},
						},
					},
				},
				gotPath: "/path/12345a/abcde1",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "12345a",
				"param2": "abcde1",
			},
		},
		{
			name: "url path template doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^\\d{5}$",
							},
							"param2": {
								DoesNotMatch: "^[a-z]{5}$",
							},
						},
					},
				},
				gotPath: "/path/12345/abcde1",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Contains: "123",
							},
							"param2": {
								Contains: "abc",
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
			name: "url path template contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								Contains: "123",
							},
							"param2": {
								Contains: "abc",
							},
						},
					},
				},
				gotPath: "/path/12345/def",
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "url path template doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "123",
							},
							"param2": {
								DoesNotContain: "abc",
							},
						},
					},
				},
				gotPath: "/path/456/def",
			},
			want: true,
			wantMap: map[string]string{
				"param1": "456",
				"param2": "def",
			},
		},
		{
			name: "url path template doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						URLPathTemplate: "/path/{param1}/{param2}",
						PathParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "123",
							},
							"param2": {
								DoesNotContain: "abc",
							},
						},
					},
				},
				gotPath: "/path/456/abcdef",
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
			name: "equalTo match",
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
			name: "equalTo does not match",
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
					"param1": []string{"12345a"},
				},
				gotQuery: url.Values{
					"param1": []string{"12345a"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "matches match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^\\d{5}$",
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
			name: "matches does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Matches: "^\\d{5}$",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"12345a"},
				},
				gotQuery: url.Values{
					"param1": []string{"12345a"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^\\d{5}$",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"12345a"},
				},
				gotQuery: url.Values{
					"param1": []string{"12345a"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "12345a",
			},
		},
		{
			name: "doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotMatch: "^\\d{5}$",
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
			name: "contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Contains: "123",
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
			name: "contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								Contains: "123",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"456"},
				},
				gotQuery: url.Values{
					"param1": []string{"456"},
				},
			},
			want:    false,
			wantMap: nil,
		},
		{
			name: "doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "123",
							},
						},
					},
				},
				gotRawQuery: url.Values{
					"param1": []string{"456"},
				},
				gotQuery: url.Values{
					"param1": []string{"456"},
				},
			},
			want: true,
			wantMap: map[string]string{
				"param1": "456",
			},
		},
		{
			name: "doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								DoesNotContain: "123",
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
			name: "contains and doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Contains:       "12345",
							DoesNotContain: "abcde",
						},
					},
				},
				gotBody: "12345",
			},
			want: true,
		},
		{
			name: "contains and doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Contains:       "12345",
							DoesNotContain: "abcde",
						},
					},
				},
				gotBody: "12345abcde",
			},
			want: false,
		},
		{
			name: "matches and doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Matches:      "^\\d{5}$",
							DoesNotMatch: "[a-zA-Z]",
						},
					},
				},
				gotBody: "12345",
			},
			want: true,
		},
		{
			name: "empty",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{},
				},
				gotBody: "12345",
			},
			want: true,
		},
		{
			name: "equalTo match",
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
		{
			name: "equalTo does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							EqualTo: "12345",
						},
					},
				},
				gotBody: "123456",
			},
			want: false,
		},
		{
			name: "matches match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Matches: "^\\d{5}$",
						},
					},
				},
				gotBody: "12345",
			},
			want: true,
		},
		{
			name: "matches does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Matches: "^\\d{5}$",
						},
					},
				},
				gotBody: "12345a",
			},
			want: false,
		},
		{
			name: "doesNotMatch match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotMatch: "^\\d{5}$",
						},
					},
				},
				gotBody: "12345a",
			},
			want: true,
		},
		{
			name: "doesNotMatch does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotMatch: "^\\d{5}$",
						},
					},
				},
				gotBody: "12345",
			},
			want: false,
		},
		{
			name: "contains match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Contains: "123",
						},
					},
				},
				gotBody: "12345",
			},
			want: true,
		},
		{
			name: "contains does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							Contains: "123",
						},
					},
				},
				gotBody: "456",
			},
			want: false,
		},
		{
			name: "doesNotContain match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotContain: "123",
						},
					},
				},
				gotBody: "456",
			},
			want: true,
		},
		{
			name: "doesNotContain does not match",
			args: args{
				endpoint: model.Endpoint{
					Request: model.Request{
						Body: model.Matcher{
							DoesNotContain: "123",
						},
					},
				},
				gotBody: "12345",
			},
			want: false,
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
