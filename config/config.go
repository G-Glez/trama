package config

import (
	"os"
)

type Config struct {
	Port        string
	GinMode     string
	DatabasePath string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		GinMode:      getEnv("GIN_MODE", "debug"),
		DatabasePath: getEnv("DATABASE_PATH", "data/trama.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
