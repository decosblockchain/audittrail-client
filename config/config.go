package config

import (
	"os"
)

func ServerUrl() string {
	return getEnv("SERVER", "http://audittrail-server:3000/")
}

func SendUrl() string {
	return ServerUrl() + "send"
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
