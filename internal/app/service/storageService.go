package service

import (
	"github.com/EMus88/GO-Yandex-Study/internal/repository"
	"github.com/EMus88/GO-Yandex-Study/pkg/idgenerator"
)

type StorageService struct {
	repos *repository.URLStorage
}

func NewStorageService(repos *repository.URLStorage) *StorageService {
	return &StorageService{repos: repos}
}

func (ss *StorageService) SaveURL(value string) string {
	key := idgenerator.CreateID()
	ss.repos.SaveURL(*key, value)
	return *key
}

func (ss *StorageService) GetURLbyID(key string) string {
	value := ss.repos.GetURLbyID(key)
	return value
}
