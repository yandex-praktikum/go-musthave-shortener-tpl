package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string  `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       url.URL `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StorageFile   string  `env:"FILE_STORAGE_PATH" envDefault:"urlStorage.gob"`
}

func Load() (*Config, error) {
	config := &Config{}

	if errEnv := env.Parse(config); errEnv != nil {
		return nil, fmt.Errorf("cannot parse config from environment: %w", errEnv)
	}

	overrideWithCliParams(config)

	return config, nil
}

func overrideWithCliParams(config *Config) {
	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	flag.StringVar(&config.StorageFile, "f", config.StorageFile, "File storage path")
	flag.Func("b", "Base URL", func(flagValue string) error {
		if flagValue == "" {
			return nil
		}

		baseURL, errParse := url.Parse(flagValue)
		if errParse != nil {
			return fmt.Errorf("cannot parse [%s] as URL: %w", flagValue, errParse)
		}
		config.BaseURL = *baseURL

		return nil
	})
	flag.Parse()
}
