package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var envFileLoaded = false

type Config struct {
	LogLevel          slog.Level `env:"LOG_LEVEL,required"`
	StorageServerAddr string     `env:"STORAGE_SERVER_ADDR,required"`
}

func New() (*Config, error) {
	if !envFileLoaded {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("godotenv: %w", err)
		}

		envFileLoaded = true
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env: %w", err)
	}

	return cfg, nil
}