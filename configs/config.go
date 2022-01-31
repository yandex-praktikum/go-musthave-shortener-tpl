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
	DBConnectionStr string
}

const Secret string = "skmsdfoiumasdfmasmdnfklwaeklasdf"

func NewConfig() *Config {

	//load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	var serverAdress string
	var baseURL string
	var fileStoragePath string
	var dbConnecionStr string

	serverAddressDefault := os.Getenv("SERVER_ADDRESS")
	baseURLDefault := os.Getenv("BASE_URL")
	fileStoragePathDefault := os.Getenv("FILE_STORAGE_PATH")
	dbConnecionStrDefault := os.Getenv("DATABASE_DSN")

	flag.StringVar(&serverAdress, "a", serverAddressDefault, "address of API server")
	flag.StringVar(&baseURL, "b", baseURLDefault, "base URL for short URL")
	flag.StringVar(&fileStoragePath, "f", fileStoragePathDefault, "path to storage file")
	flag.StringVar(&dbConnecionStr, "d", dbConnecionStrDefault, "str to DB connection")

	flag.Parse()

	return &Config{
		ServerAdress:    serverAdress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		DBConnectionStr: dbConnecionStr,
	}
}

func NewConfigForTest() *Config {
	return &Config{
		ServerAdress:    "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "urls.data",
		DBConnectionStr: "postgres://user1:123@localhost:5432/mydb1",
	}
}
