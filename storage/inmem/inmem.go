package inmem

import (
	"sync"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type InMemStorage struct {
	sync.RWMutex

	store []model.StoreURL
}

func New() *InMemStorage {
	return &InMemStorage{
		RWMutex: sync.RWMutex{},
		store:   make([]model.StoreURL, 0),
	}
}

func (s *InMemStorage) GetByID(id int) *model.StoreURL {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	for _, url := range s.store {
		if url.ID == id {
			return &url
		}
	}

	return nil
}

func (s *InMemStorage) Save(u model.StorableURL) model.StoreURL {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	id := len(s.store)
	storedURL := model.StoreURL{
		ID:      id,
		LongURL: u.LongURL,
	}
	s.store = append(s.store, storedURL)

	return storedURL
}

func (s *InMemStorage) GetAll() []model.StoreURL {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.store
}

func (s *InMemStorage) Load(u model.StoreURL) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.store = append(s.store, u)
}
