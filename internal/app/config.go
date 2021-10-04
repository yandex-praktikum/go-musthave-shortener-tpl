package app

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

func LoadConfig() (Config, error) {
	var conf Config
	if errConf := env.Parse(&conf); errConf != nil {
		return conf, fmt.Errorf("cannot parse config from environment: %w", errConf)
	}

	var baseURLStr string
	flag.StringVar(&conf.ServerAddress, "a", conf.ServerAddress, "Server address")
	flag.StringVar(&conf.StorageFile, "f", conf.StorageFile, "File storage path")
	flag.StringVar(&baseURLStr, "b", conf.BaseURL.String(), "Base URL")
	flag.Parse()

	if baseURLStr > "" {
		baseURL, errParse := url.Parse(baseURLStr)
		if errParse != nil {
			return conf, fmt.Errorf("cannot parse base URL: %w", errParse)
		}
		conf.BaseURL = *baseURL
	}

	return conf, nil
}
