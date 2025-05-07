package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	main "github.com/dev-shimada/gostubby"
)

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
