package v1

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type Service struct {
	Storage storage.ShortenerStorage
	BaseURL url.URL
}

func New(s storage.ShortenerStorage, u url.URL) (*Service, error) {
	if s == nil {
		return nil, errors.New("storage should not be nil")
	}
	return &Service{s, u}, nil
}

func (s *Service) ShortenURL(u model.URLToShorten) (*model.ShortenedURL, error) {
	url, err := s.Storage.SaveURL(u)
	if err != nil {
		return nil, fmt.Errorf("cannot shorten url: %w", err)
	}
	log.Printf("Shortened: %s", url)

	return &url, nil
}

func (s *Service) GetURLByID(id int) (*model.ShortenedURL, error) {
	return s.Storage.GetURLByID(id)
}

func (s *Service) LookupURL(u url.URL) (*model.ShortenedURL, error) {
	return s.Storage.LookupURL(u)
}

func (s *Service) GetUserURLs(userID int64) ([]model.ShortenedURL, error) {
	return s.Storage.ListByUserID(userID)
}

func (s *Service) AbsoluteURL(u model.ShortenedURL) (*url.URL, error) {
	urlPath := fmt.Sprintf("%d", u.ID)

	shortURL, err := s.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, fmt.Errorf("cannot make absolute URL for id [%d]", u.ID)
	}

	return shortURL, nil
}
