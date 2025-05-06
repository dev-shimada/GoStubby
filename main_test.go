package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	main "github.com/dev-shimada/gostubby"
	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/dev-shimada/gostubby/internal/infrastructure/config"
	"github.com/google/go-cmp/cmp"
)

func Test_pathMatcher(t *testing.T) {
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
			name: "url",
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
		endpoint model.Endpoint
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
				endpoint: model.Endpoint{
					Request: model.Request{
						QueryParameters: map[string]model.Matcher{
							"param1": {
								EqualTo: "12345",
							},
						},
					},
				},
				gotQuery: url.Values{
					"param1": []string{"12345"},
				},
			},
			want: true,
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

func Test_debugPathMatcherDetailed(t *testing.T) {
	// Load the configuration
	repo := config.NewConfigRepository()
	endpoints, err := repo.Load("testdata/test_config.json")
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

	// Check each parameter constraint
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

func Test_handle(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "not found",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/not-found", nil),
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			main.Handle(w, tt.args.r)

			resp := w.Result()
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			}()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("handle() status = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
