package helpers

import (
	"net/http"
	"slices"
	"testing"

	"github.com/tanmancan/draw-together/internal/config"
)

func TestCheckOriginHandler(t *testing.T) {
	notAllowedOrigin := "https://not-allowed"
	allowedOrigin := "https://allowed-origin"

	config.AppConfig.CorsConfig.AllowedOrigin = append(config.AppConfig.CorsConfig.AllowedOrigin, allowedOrigin)
	slices.Sort[[]string](config.AppConfig.CorsConfig.AllowedOrigin)
	testUrl := "localhost/test"

	reqNotAllowed, _ := http.NewRequest("GET", testUrl, nil)
	reqNotAllowed.Header.Set("Origin", notAllowedOrigin)

	reqAllowed, _ := http.NewRequest("POST", testUrl, nil)
	reqAllowed.Header.Set("Origin", allowedOrigin)

	reqNoOrigin, _ := http.NewRequest("PUT", testUrl, nil)

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "CheckOriginHandler returns false if request origin is NOT in approved list",
			args: args{
				r: reqNotAllowed,
			},
			want: false,
		},
		{
			name: "CheckOriginHandler returns true if request origin is in approved list",
			args: args{
				r: reqAllowed,
			},
			want: true,
		},
		{
			name: "CheckOriginHandler returns false if request does not contain origin header",
			args: args{
				r: reqNoOrigin,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckOriginHandler(tt.args.r); got != tt.want {
				t.Errorf("CheckOriginHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
