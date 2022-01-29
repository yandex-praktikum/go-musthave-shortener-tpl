package service

import (
	"fmt"
	"log"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/EMus88/go-musthave-shortener-tpl/pkg/idgenerator"
)

type Repository interface {
	SaveURL(shortModel *model.Shorten, sessionID string) error
	GetURL(id string) string
	PingDB() error
	GetCookie(s string) bool
	SaveCookie(s string) error
}
type Service struct {
	Repository
	Config configs.Config
	Auth   Auth
}

func NewService(repos *repository.Storage, config *configs.Config) *Service {
	return &Service{Repository: repos, Config: *config}
}

//save long URL in stotage and return short URL
func (s *Service) SaveURL(longURL string, sessionID string) (string, error) {
	var shortModel model.Shorten
	shortModel.URLID = idgenerator.CreateID(8)
	shortModel.ShortURL = fmt.Sprint(s.Config.BaseURL, "/", shortModel.URLID)
	shortModel.LongURL = longURL

	key, _ := s.Auth.ReadSessionID(sessionID)

	if err := s.Repository.SaveURL(&shortModel, key); err != nil {
		return "", err
	}

	//save to file
	// file, err := os.OpenFile(s.Config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	// if err != nil {
	// 	return "", err
	// }
	// defer file.Close()

	// s.Model.ID = key
	// s.Model.LongURL = value

	// data, err := json.MarshalIndent(s.Model, "", " ")
	// if err != nil {
	// 	return "", err
	// }
	// file.Write(data)

	return shortModel.ShortURL, nil
}

//get long URL from stotage by short URL
func (s *Service) GetURL(key string) (string, error) {
	originURL := s.Repository.GetURL(key)
	return originURL, nil
}

func (s *Service) CreateNewSession() (string, error) {
	id, encID, err := s.Auth.CreateSissionID()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Repository.SaveCookie(id); err != nil {
		return "", err
	}

	return encID, nil
}

// func (s *Service) LoadFromFile() {

// 	var model model.File
// 	file, err := os.ReadFile(s.Config.FileStoragePath)
// 	if err != nil {
// 		return
// 	}
// 	str := strings.Split(string(file), "}")
// 	for i := 0; i < (len(str) - 1); i++ {

// 		if err := json.Unmarshal([]byte(str[i]+"}"), &model); err != nil {
// 			log.Fatal(err)
// 		}
// 		s.Repository.SaveURL(model.ID, model.LongURL)
// 	}
// }
