package inmem

import (
	"sync"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type InMemStorage struct {
	sync.RWMutex

	store []model.ShortenedURL
}

func New() *InMemStorage {
	return &InMemStorage{
		RWMutex: sync.RWMutex{},
		store:   make([]model.ShortenedURL, 0),
	}
}

func (s *InMemStorage) GetByID(id int) *model.ShortenedURL {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	for _, url := range s.store {
		if url.ID == id {
			return &url
		}
	}

	return nil
}

func (s *InMemStorage) ListByUserID(userID int) []model.ShortenedURL {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	result := make([]model.ShortenedURL, 0)
	for _, url := range s.store {
		if url.UserID == userID {
			result = append(result, url)
		}
	}

	return result
}

func (s *InMemStorage) Save(u model.URLToShorten) model.ShortenedURL {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	id := len(s.store)
	storedURL := model.ShortenedURL{
		UserID:  u.UserID,
		ID:      id,
		LongURL: u.LongURL,
	}
	s.store = append(s.store, storedURL)

	return storedURL
}

func (s *InMemStorage) GetAll() []model.ShortenedURL {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.store
}

func (s *InMemStorage) Load(u model.ShortenedURL) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.store = append(s.store, u)
}
