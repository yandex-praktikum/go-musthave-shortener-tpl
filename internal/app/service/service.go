package service

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/models/file"
	"github.com/EMus88/go-musthave-shortener-tpl/pkg/idgenerator"
)

type Repository interface {
	SaveURL(key string, value string)
	GetURL(id string) (string, error)
}
type Service struct {
	Repository
	Model  file.Model
	Config configs.Config
}

func NewService(repos *repository.URLStorage, model *file.Model, config *configs.Config) *Service {
	return &Service{Repository: repos, Model: *model, Config: *config}
}

//save long URL in stotage and return short URL
func (s *Service) SaveURL(value string) (string, error) {
	//save to map
	key := idgenerator.CreateID()
	s.Repository.SaveURL(key, value)
	//save to file
	file, err := os.OpenFile(s.Config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return "", err
	}
	defer file.Close()
	s.Model.ID = key
	s.Model.LongURL = value

	data, err := json.MarshalIndent(s.Model, "", " ")
	if err != nil {
		return "", err
	}
	file.Write(data)

	return key, nil
}

//get long URL from stotage by short URL
func (s *Service) GetURL(key string) (string, error) {
	value, err := s.Repository.GetURL(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *Service) LoadFromFile() {

	var model file.Model
	file, err := os.ReadFile(s.Config.FileStoragePath)
	if err != nil {
		return
	}
	str := strings.Split(string(file), "}")
	for i := 0; i < (len(str) - 1); i++ {

		if err := json.Unmarshal([]byte(str[i]+"}"), &model); err != nil {
			log.Fatal(err)
		}
		s.Repository.SaveURL(model.ID, model.LongURL)
	}
}
