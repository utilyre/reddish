package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Mode Mode `env:"-"`

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

	em := make(map[string]string)
	if mode == ModeDev {
		var err error
		em, err = godotenv.Read(".env")
		if err != nil {
			return nil, fmt.Errorf("godotenv: %w", err)
		}
	}

	cfg, err := env.ParseAsWithOptions[Config](env.Options{Environment: em})
	if err != nil {
		return nil, fmt.Errorf("env: %w", err)
	}

	cfg.Mode = mode
	fmt.Printf("Config: %+v\n", cfg)

	return &cfg, nil
}
