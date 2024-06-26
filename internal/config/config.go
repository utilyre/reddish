package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Mode     Mode       `env:"-"`
	LogLevel slog.Level `env:"LOG_LEVEL,required"`

	StorageServerAddr string `env:"STORAGE_SERVER_ADDR,required"`
}

func New() (*Config, error) {
	modeStr, ok := os.LookupEnv("MODE")
	if !ok {
		modeStr = "DEV"
	}

	var mode Mode
	if err := mode.UnmarshalText([]byte(modeStr)); err != nil {
		return nil, err
	}

	if mode == ModeDev {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("godotenv: %w", err)
		}
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("env: %w", err)
	}

	cfg.Mode = mode

	return &cfg, nil
}
