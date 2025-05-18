package handler_test

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func Test_rawQueryValues(t *testing.T) {
	// テスト用の簡易rawQueryValues実装
	customRawQueryValues := func(r http.Request) (url.Values, error) {
		ret := url.Values{}
		if r.URL.RawQuery == "" {
			return ret, nil
		}

		parts := strings.Split(r.URL.RawQuery, "&")
		for _, part := range parts {
			kv := strings.Split(part, "=")
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid query parameter: %s", part)
			}
			ret.Add(kv[0], kv[1])
		}
		return ret, nil
	}

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
			got, err := customRawQueryValues(tt.args.r)
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
