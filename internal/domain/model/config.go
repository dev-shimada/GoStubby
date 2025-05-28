package model

import (
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

// define the structure of the JSON configuration file
type Matcher struct {
	EqualTo        any `json:"equalTo"`
	Matches        any `json:"matches"`
	DoesNotMatch   any `json:"doesNotMatch"`
	Contains       any `json:"contains"`
	DoesNotContain any `json:"doesNotContain"`
}
type Request struct {
	URL             string `json:"url"`             // パスパラメータ、クエリパラメータを含む完全一致
	URLPattern      string `json:"urlPattern"`      // パスパラメータ、クエリパラメータを含む正規表現での完全一致
	URLPath         string `json:"urlPath"`         // パスパラメータを含む完全一致
	URLPathPattern  string `json:"urlPathPattern"`  // パスパラメータを含む正規表現での完全一致
	URLPathTemplate string `json:"urlPathTemplate"` // パスパラメータを含むテンプレートでの完全一致

	Method          string             `json:"method"`
	Headers         map[string]Matcher `json:"headers"` // HTTP header matchers
	QueryParameters map[string]Matcher `json:"queryParameters"`
	PathParameters  map[string]Matcher `json:"pathParameters"`
	Body            Matcher            `json:"body"`
}
type Response struct {
	Status        int               `json:"status"`
	BodyFileName  string            `json:"bodyFileName"` // bodyFileNameが指定されている場合は、bodyは無視される
	Body          string            `json:"body"`         // bodyFileNameが指定されていない場合は、bodyを使用する
	Headers       map[string]string `json:"headers"`
	Transformaers []string          `json:"transformers"`
}
type Endpoint struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
}

func (endpoint Endpoint) PathMatcher(gotRawPath, gotPath string) (bool, map[string]string) {
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
		switch {
		case v.EqualTo != nil && gotPathUnits[posMap[k]] != fmt.Sprint(v.EqualTo):
			return false, nil
		case v.Matches != nil && !regexp.MustCompile(v.Matches.(string)).MatchString(gotPathUnits[posMap[k]]):
			return false, nil
		case v.DoesNotMatch != nil && regexp.MustCompile(v.DoesNotMatch.(string)).MatchString(gotPathUnits[posMap[k]]):
			return false, nil
		case v.Contains != nil && !strings.Contains(gotPathUnits[posMap[k]], v.Contains.(string)):
			return false, nil
		case v.DoesNotContain != nil && strings.Contains(gotPathUnits[posMap[k]], v.DoesNotContain.(string)):
			return false, nil
		}
	}
	ret := make(map[string]string)
	for k, v := range posMap {
		ret[k] = gotPathUnits[v]
	}
	return true, ret
}

func (endpoint Endpoint) QueryMatcher(gotRawQuery, gotQuery url.Values) (bool, map[string]string) {
	for k, v := range endpoint.Request.QueryParameters {
		switch {
		case v.EqualTo != nil && gotRawQuery.Get(k) != fmt.Sprint(v.EqualTo):
			return false, nil
		case v.Matches != nil && !regexp.MustCompile(v.Matches.(string)).MatchString(gotRawQuery.Get(k)):
			return false, nil
		case v.DoesNotMatch != nil && regexp.MustCompile(v.DoesNotMatch.(string)).MatchString(gotRawQuery.Get(k)):
			return false, nil
		case v.Contains != nil && !strings.Contains(gotRawQuery.Get(k), v.Contains.(string)):
			return false, nil
		case v.DoesNotContain != nil && strings.Contains(gotRawQuery.Get(k), v.DoesNotContain.(string)):
			return false, nil
		}
	}
	ret := make(map[string]string)
	for k, v := range gotQuery {
		ret[k] = v[0]
	}
	return true, ret
}

func (endpoint Endpoint) BodyMatcher(body string) bool {
	switch {
	case endpoint.Request.Body.EqualTo != nil && body != fmt.Sprint(endpoint.Request.Body.EqualTo):
		return false
	case endpoint.Request.Body.Matches != nil && !regexp.MustCompile(endpoint.Request.Body.Matches.(string)).MatchString(body):
		return false
	case endpoint.Request.Body.DoesNotMatch != nil && regexp.MustCompile(endpoint.Request.Body.DoesNotMatch.(string)).MatchString(body):
		return false
	case endpoint.Request.Body.Contains != nil && !strings.Contains(body, endpoint.Request.Body.Contains.(string)):
		return false
	case endpoint.Request.Body.DoesNotContain != nil && strings.Contains(body, endpoint.Request.Body.DoesNotContain.(string)):
		return false
	}
	return true
}

func (endpoint Endpoint) HeaderMatcher(headers map[string][]string) bool {
	for k, v := range endpoint.Request.Headers {
		headerVal := ""
		if values, exists := headers[k]; exists && len(values) > 0 {
			headerVal = values[0]
		}

		switch {
		case v.EqualTo != nil && headerVal != fmt.Sprint(v.EqualTo):
			return false
		case v.Matches != nil && !regexp.MustCompile(v.Matches.(string)).MatchString(headerVal):
			return false
		case v.DoesNotMatch != nil && regexp.MustCompile(v.DoesNotMatch.(string)).MatchString(headerVal):
			return false
		case v.Contains != nil && !strings.Contains(headerVal, v.Contains.(string)):
			return false
		case v.DoesNotContain != nil && strings.Contains(headerVal, v.DoesNotContain.(string)):
			return false
		}
	}
	return true
}
