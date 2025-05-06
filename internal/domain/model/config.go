package model

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
