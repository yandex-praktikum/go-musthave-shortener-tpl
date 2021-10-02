package app

import (
	"net/url"
	"sync"
)

type Repository interface {
	GetURLBy(id int) *url.URL
	SaveURL(u url.URL) int
}

type MemRepository struct {
	sync.RWMutex

	store map[int]url.URL
}

func NewMemRepository() Repository {
	return &MemRepository{
		RWMutex: sync.RWMutex{},
		store:   make(map[int]url.URL),
	}
}

func (r *MemRepository) SaveURL(u url.URL) int {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	id := len(r.store)
	r.store[id] = u

	return id
}

func (r *MemRepository) GetURLBy(id int) *url.URL {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	longURL, ok := r.store[id]
	if !ok {
		return nil
	}
	return &longURL
}
