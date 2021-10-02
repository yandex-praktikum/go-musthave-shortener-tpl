package app

import (
	"flag"
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

	var baseURLStr string
	flag.StringVar(&conf.ServerAddress, "a", conf.ServerAddress, "Server address")
	flag.StringVar(&conf.StorageFile, "f", conf.StorageFile, "File storage path")
	flag.StringVar(&baseURLStr, "b", conf.BaseURL.String(), "Base URL")
	flag.Parse()

	if baseURLStr > "" {
		baseURL, errParse := url.Parse(baseURLStr)
		if errParse != nil {
			panic(errParse)
		}
		conf.BaseURL = *baseURL
	}

	return conf
}
