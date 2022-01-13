package configs

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAdress    string
	BaseURL         string
	FileStoragePath string
}

func NewConfig() *Config {

	//load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	var serverAdress string
	var baseURL string
	var fileStoragePath string

	serverAddressDefault := os.Getenv("SERVER_ADDRESS")
	baseURLDefault := os.Getenv("BASE_URL")
	fileStoragePathDefault := os.Getenv("FILE_STORAGE_PATH")

	flag.StringVar(&serverAdress, "a", serverAddressDefault, "address of API server")
	flag.StringVar(&baseURL, "b", baseURLDefault, "base URL for short URL")
	flag.StringVar(&fileStoragePath, "f", fileStoragePathDefault, "path to storage file")
	flag.Parse()

	return &Config{
		ServerAdress:    serverAdress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
	}
}
func NewConfigForTest() *Config {
	return &Config{
		ServerAdress:    "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "urls.data",
	}
}
