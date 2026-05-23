package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseDSN  string
	AuthUsername string
	AuthPassword string
	Listen       string
	Env          string
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseDSN:   strings.TrimSpace(os.Getenv("DATABASE_DSN")),
		AuthUsername:  strings.TrimSpace(os.Getenv("AUTH_USERNAME")),
		AuthPassword:  os.Getenv("AUTH_PASSWORD"),
		Listen: envOr("APP_LISTEN", ":8080"),
		Env:    envOr("APP_ENV", "production"),
	}
	if cfg.DatabaseDSN == "" {
		return cfg, fmt.Errorf("DATABASE_DSN is required")
	}
	if cfg.AuthUsername == "" || cfg.AuthPassword == "" {
		return cfg, fmt.Errorf("AUTH_USERNAME and AUTH_PASSWORD are required")
	}
	return cfg, nil
}

func (c Config) Development() bool {
	return c.Env == "development"
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
