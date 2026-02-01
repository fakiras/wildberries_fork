package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort int
	DSN      string
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("HTTP_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://localhost:5432/wildberries?sslmode=disable"
	}
	return &Config{
		HTTPPort: port,
		DSN:      dsn,
	}
}
