package handler_test

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/dev-shimada/gostubby/internal/handler"
)

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
