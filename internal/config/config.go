package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var envFileLoaded = false

type Config struct {
	StorageServerAddr string `env:"STORAGE_SERVER_ADDR,required"`
}

func New() (*Config, error) {
	if !envFileLoaded && os.Getenv("MODE") != "prod" {
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
