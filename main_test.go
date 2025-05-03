package main_test

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	main "github.com/dev-shimada/gostubby"
	"github.com/google/go-cmp/cmp"
)

func Test_pathMatcher(t *testing.T) {
	type args struct {
		endpoint   main.Endpoint
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
			name: "url",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URL: "http://example.com/path",
					},
				},
				gotRawPath: "http://example.com/path/",
			},
			want: true,
		},
		{
			name: "url false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URL: "http://example.com/path",
					},
				},
				gotRawPath: "http://example.com/path?a=1",
			},
			want: false,
		},
		{
			name: "urlPattern",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPattern: "http://example.com/(\\d{5})/",
					},
				},
				gotRawPath: "http://example.com/123456",
			},
			want: true,
		},
		{
			name: "urlPattern false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPattern: "http://example.com/(\\d{5})/",
					},
				},
				gotRawPath: "http://example.com/abcde",
			},
			want: false,
		},
		{
			name: "urlPath",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPath: "http://example.com/path",
					},
				},
				gotPath: "http://example.com/path",
			},
			want: true,
		},
		{
			name: "urlPath false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPath: "http://example.com/path",
					},
				},
				gotPath: "http://example.com/path/abc",
			},
			want: false,
		},
		{
			name: "urlPathPattern",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathPattern: "http://example.com/(\\d{5})",
					},
				},
				gotPath: "http://example.com/12345/abc",
			},
			want: true,
		},
		{
			name: "urlPathPattern false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathPattern: "http://example.com/(\\d{5})",
					},
				},
				gotPath: "http://example.com/abc",
			},
			want: false,
		},
		{
			name: "urlPathTemplate",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathTemplate: "http://example.com/{path1}/{path2}/{path3}/{path4}/{path5}",
						PathParameters: map[string]main.Matcher{
							"path1": {
								EqualTo: "12345",
							},
							"path2": {
								Matches: "[0-9]{5}",
							},
							"path3": {
								DoesNotMatch: "[a-z]{5}",
							},
							"path4": {
								Contains: "abc",
							},
							"path5": {
								DoesNotContain: "xyz",
							},
						},
					},
				},
				gotPath: "http://example.com/12345/67890/00000/abcxyz/12345",
			},
			want: true,
			wantMap: map[string]string{
				"path1": "12345",
				"path2": "67890",
				"path3": "00000",
				"path4": "abcxyz",
				"path5": "12345",
			},
		},
		{
			name: "urlPathTemplate equalTo false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathTemplate: "http://example.com/{path}",
						PathParameters: map[string]main.Matcher{
							"path": {
								EqualTo: "12345",
							},
						},
					},
				},
				gotPath: "http://example.com/123456",
			},
			want: false,
		},
		{
			name: "urlPathTemplate contains false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathTemplate: "http://example.com/{path}",
						PathParameters: map[string]main.Matcher{
							"path": {
								Contains: "12345",
							},
						},
					},
				},
				gotPath: "http://example.com/1234",
			},
			want: false,
		},
		{
			name: "urlPathTemplate doesNotContain false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathTemplate: "http://example.com/{path}",
						PathParameters: map[string]main.Matcher{
							"path": {
								DoesNotContain: "12345",
							},
						},
					},
				},
				gotPath: "http://example.com/12345abc",
			},
			want: false,
		},
		{
			name: "urlPathTemplate matches false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathTemplate: "http://example.com/{path}",
						PathParameters: map[string]main.Matcher{
							"path": {
								Matches: "[0-9]{5}",
							},
						},
					},
				},
				gotPath: "http://example.com/1234",
			},
			want: false,
		},
		{
			name: "urlPathTemplate doesNotMatch false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						URLPathTemplate: "http://example.com/{path}",
						PathParameters: map[string]main.Matcher{
							"path": {
								DoesNotMatch: "[0-9]{5}",
							},
						},
					},
				},
				gotPath: "http://example.com/12345abc",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMap := main.ExportPathMatcher(tt.args.endpoint, tt.args.gotRawPath, tt.args.gotPath); got != tt.want {
				t.Errorf("pathMatcher() = %v, want %v", got, tt.want)
			} else if !cmp.Equal(gotMap, tt.wantMap) {
				t.Errorf("diff: %v", cmp.Diff(gotMap, tt.wantMap))
			}
		})
	}
}

func Test_queryMatcher(t *testing.T) {
	type args struct {
		endpoint main.Endpoint
		gotQuery url.Values
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						QueryParameters: map[string]main.Matcher{
							"param1": {
								EqualTo: "12345",
							},
							"param2": {
								Matches: "[0-9]{5}",
							},
							"param3": {
								DoesNotMatch: "[a-z]{5}",
							},
							"param4": {
								Contains: "abc",
							},
							"param5": {
								DoesNotContain: "xyz",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param1": []string{"12345"},
					"param2": []string{"67890"},
					"param3": []string{"00000"},
					"param4": []string{"abcxyz"},
					"param5": []string{"12345"},
				},
			},
			want: true,
		},
		{
			name: "equalTo false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						QueryParameters: map[string]main.Matcher{
							"param": {
								EqualTo: "12345",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param": []string{"123456"},
				},
			},
			want: false,
		},
		{
			name: "contains false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						QueryParameters: map[string]main.Matcher{
							"param": {
								Contains: "12345",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param": []string{"1234"},
				},
			},
			want: false,
		},
		{
			name: "doesNotContain false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						QueryParameters: map[string]main.Matcher{
							"param": {
								DoesNotContain: "12345",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param": []string{"12345abc"},
				},
			},
			want: false,
		},
		{
			name: "matches false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						QueryParameters: map[string]main.Matcher{
							"param": {
								Matches: "[0-9]{5}",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param": []string{"1234"},
				},
			},
			want: false,
		},
		{
			name: "doesNotMatch false",
			args: args{
				endpoint: main.Endpoint{
					Request: main.Request{
						QueryParameters: map[string]main.Matcher{
							"param": {
								DoesNotMatch: "[0-9]{5}",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param": []string{"12345abc"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := main.ExportQueryMatcher(tt.args.endpoint, tt.args.gotQuery); got != tt.want {
				t.Errorf("queryMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []main.Endpoint
		wantErr bool
	}{
		{
			name: "valid single file",
			args: args{
				path: "testdata/test_config.json",
			},
			want: []main.Endpoint{
				{
					Request: main.Request{
						URLPathTemplate: "/example/{path1}/{path2}/{path3}/{path4}/{path5}",
						Method:          "GET",
						PathParameters: map[string]main.Matcher{
							"path1": {
								EqualTo: "v1",
							},
							"path2": {
								Matches: "^[a-zA-Z0-9]{3}$",
							},
							"path3": {
								DoesNotMatch: "[a-zA-Z]{3}",
							},
							"path4": {
								Contains: "a",
							},
							"path5": {
								DoesNotContain: "b",
							},
						},
						QueryParameters: map[string]main.Matcher{
							"param1": {
								EqualTo: "value1",
							},
							"param2": {
								Matches: "[a-zA-Z0-9]{3}",
							},
							"param3": {
								DoesNotMatch: "[a-zA-Z]{3}",
							},
							"param4": {
								Contains: "a",
							},
							"param5": {
								DoesNotContain: "b",
							},
						},
					},
					Response: main.Response{
						Status: 200,
						Body:   `{"message": "This is a stub response", "param1"="{{.Query.param1}}", "param2"="{{.Query.param2}}", "param3"="{{.Query.param3}}", "param4"="{{.Query.param4}}", "param5"="{{.Query.param5}}"}` + "\n",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid directory",
			args: args{
				path: "testdata",
			},
			want: []main.Endpoint{
				{
					Request: main.Request{
						URLPathTemplate: "/example/{path1}/{path2}/{path3}/{path4}/{path5}",
						Method:          "GET",
						PathParameters: map[string]main.Matcher{
							"path1": {
								EqualTo: "v1",
							},
							"path2": {
								Matches: "^[a-zA-Z0-9]{3}$",
							},
							"path3": {
								DoesNotMatch: "[a-zA-Z]{3}",
							},
							"path4": {
								Contains: "a",
							},
							"path5": {
								DoesNotContain: "b",
							},
						},
						QueryParameters: map[string]main.Matcher{
							"param1": {
								EqualTo: "value1",
							},
							"param2": {
								Matches: "[a-zA-Z0-9]{3}",
							},
							"param3": {
								DoesNotMatch: "[a-zA-Z]{3}",
							},
							"param4": {
								Contains: "a",
							},
							"param5": {
								DoesNotContain: "b",
							},
						},
					},
					Response: main.Response{
						Status: 200,
						Body:   `{"message": "This is a stub response", "param1"="{{.Query.param1}}", "param2"="{{.Query.param2}}", "param3"="{{.Query.param3}}", "param4"="{{.Query.param4}}", "param5"="{{.Query.param5}}"}` + "\n",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid file path",
			args: args{
				path: "testdata/nonexistent.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non-json file",
			args: args{
				path: "testdata/test.txt",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := main.ExportLoadConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_debugPathMatcherDetailed(t *testing.T) {
	// Load the configuration
	endpoints, err := main.ExportLoadConfig("testdata/test_config.json")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(endpoints) == 0 {
		t.Fatalf("No endpoints loaded from configuration")
	}

	endpoint := endpoints[0]
	testPath := "/example/v1/123/000/axyz"

	// Split the paths for comparison
	pathTemplate := endpoint.Request.URLPathTemplate
	t.Logf("Template: %s", pathTemplate)
	t.Logf("Test path: %s", testPath)

	templateUnits := strings.Split(strings.TrimRight(pathTemplate, "/"), "/")
	testPathUnits := strings.Split(strings.TrimRight(testPath, "/"), "/")

	t.Logf("Template units (%d): %v", len(templateUnits), templateUnits)
	t.Logf("Test path units (%d): %v", len(testPathUnits), testPathUnits)

	// Check lengths match
	if len(templateUnits) != len(testPathUnits) {
		t.Logf("Length mismatch: template=%d, path=%d", len(templateUnits), len(testPathUnits))
	} else {
		t.Logf("Length match: %d units", len(templateUnits))
	}

	// Check path parameters
	posMap := make(map[string]int)
	for k := range endpoint.Request.PathParameters {
		placeHolder := fmt.Sprintf("{%s}", k)
		for i, unit := range templateUnits {
			if unit == placeHolder {
				posMap[k] = i
				if i < len(testPathUnits) {
					t.Logf("Parameter %s at position %d = %s", k, i, testPathUnits[i])
				} else {
					t.Logf("Parameter %s at position %d is out of range for test path", k, i)
				}
				break
			}
		}
	}

	t.Logf("Position map: %v", posMap)

	// Check each parameter constraint - safely this time
	for k, v := range endpoint.Request.PathParameters {
		pos, ok := posMap[k]
		if !ok {
			t.Logf("Path parameter %s not found in template", k)
			continue
		}

		if pos >= len(testPathUnits) {
			t.Logf("Position %d for parameter %s is out of range for test path", pos, k)
			continue
		}

		paramValue := testPathUnits[pos]
		t.Logf("Checking parameter %s = %s", k, paramValue)

		if v.EqualTo != nil {
			isEqual := paramValue == fmt.Sprint(v.EqualTo)
			t.Logf("  EqualTo %v: %v", v.EqualTo, isEqual)
		}
		if v.Matches != nil {
			doesMatch := regexp.MustCompile(v.Matches.(string)).MatchString(paramValue)
			t.Logf("  Matches %v: %v", v.Matches, doesMatch)
		}
		if v.DoesNotMatch != nil {
			doesNotMatch := !regexp.MustCompile(v.DoesNotMatch.(string)).MatchString(paramValue)
			t.Logf("  DoesNotMatch %v: %v", v.DoesNotMatch, doesNotMatch)
		}
		if v.Contains != nil {
			contains := strings.Contains(paramValue, v.Contains.(string))
			t.Logf("  Contains %v: %v", v.Contains, contains)
		}
		if v.DoesNotContain != nil {
			doesNotContain := !strings.Contains(paramValue, v.DoesNotContain.(string))
			t.Logf("  DoesNotContain %v: %v", v.DoesNotContain, doesNotContain)
		}
	}

	// Just directly test the path matcher function with our test path
	matched, pathMap := main.ExportPathMatcher(endpoint, "", testPath)
	t.Logf("Path matcher result: %v, map: %v", matched, pathMap)
}
