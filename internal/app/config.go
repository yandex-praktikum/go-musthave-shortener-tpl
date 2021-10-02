package app

import (
	"net/url"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string  `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       url.URL `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StorageFile   string  `env:"FILE_STORAGE_PATH" envDefault:"urlStorage.gob"`
}

func LoadConfig() Config {
	var conf Config
	errConf := env.Parse(&conf)
	if errConf != nil {
		panic(errConf)
	}
	return conf
}
