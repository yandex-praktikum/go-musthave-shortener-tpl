package service

import (
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/EMus88/go-musthave-shortener-tpl/pkg/idgenerator"
)

type Repositiries interface {
	SaveURLtoStorage(key string, value string)
	GetURLfromStorage(id string) string
}
type Service struct {
	Repositiries
}

func NewService(repos *repository.URLStorage) *Service {
	return &Service{Repositiries: repos}
}

func (s *Service) SaveURL(value string) string {
	key := idgenerator.CreateID()
	s.Repositiries.SaveURLtoStorage(key, value)
	return key
}

func (s *Service) GetURL(key string) string {
	value := s.Repositiries.GetURLfromStorage(key)
	return value
}
