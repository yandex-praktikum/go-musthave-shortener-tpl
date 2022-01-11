package repository

import (
	"errors"
	"sync"
)

type URLStorage struct {
	storage sync.Map
}

func NewStorage() *URLStorage {
	return &URLStorage{
		storage: sync.Map{},
	}
}
func (us *URLStorage) SaveURL(key string, value string) {
	us.storage.Store(key, value)
}

func (us *URLStorage) GetURL(key string) (string, error) {
	value, ok := us.storage.Load(key)
	if ok {
		return value.(string), nil
	}
	return "", errors.New("URL not found in base")
}
