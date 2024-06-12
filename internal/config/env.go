package config

import (
	"os"
	"strconv"
)

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}

	val, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}

	return val
}

func getEnvBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}

	val, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}

	return val
}

func getEnvStr(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}

	return v
}
