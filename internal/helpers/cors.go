package helpers

import (
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
)

var loggerCors = basiclogger.BasicLogger{Namespace: "internal.helpers.cors"}

func init() {
	slices.Sort[[]string](config.AppConfig.CorsConfig.AllowedOrigin)
}

func CheckOriginHandler(r *http.Request) bool {
	reqID := uuid.NewString()
	origin := r.Header.Get("Origin")

	if origin == "" {
		loggerCors.LogWarn(reqID, "origin header not provided", "origin", origin, "found", false)
		return false
	}

	allowedOrigin := config.AppConfig.CorsConfig.AllowedOrigin
	loggerCors.LogInfo(reqID, "check origin start", "allowedOrigin", allowedOrigin)
	_, found := slices.BinarySearch[[]string](allowedOrigin, origin)
	loggerCors.LogInfo(reqID, "check origin done", "allowedOrigin", allowedOrigin, "found", found)

	return found
}
