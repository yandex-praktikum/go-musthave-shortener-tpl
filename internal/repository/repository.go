package repository

import "sync"

type URLStorage struct {
	storage map[string]string
	mx      sync.Mutex
}

func NewStorage() *URLStorage {
	return &URLStorage{
		storage: make(map[string]string, 10),
		mx:      sync.Mutex{},
	}
}
func (us *URLStorage) SaveURL(key string, value string) {
	us.mx.Lock()
	us.storage[key] = value
	us.mx.Unlock()
}

func (us *URLStorage) GetURLbyID(key string) string {
	return us.storage[key]
}
