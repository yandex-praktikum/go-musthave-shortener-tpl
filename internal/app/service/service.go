package service

import "github.com/EMus88/GO-Yandex-Study/internal/repository"

type UseStorage interface {
	SaveURL(value string) string
	GetURLbyID(id string) string
}
type Service struct {
	UseStorage
}

func NewService(repos *repository.URLStorage) *Service {
	return &Service{UseStorage: NewStorageService(repos)}
}
