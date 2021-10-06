package v1

import (
	"fmt"
	"log"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type Service struct {
	Storage storage.Storage
	BaseURL url.URL
}

func New(s storage.Storage, u url.URL) *Service {
	return &Service{s, u}
}

func (s *Service) ShortenURL(newURL model.StorableURL) (*url.URL, error) {
	url := s.Storage.Save(newURL)
	urlPath := fmt.Sprintf("%d", url.ID)

	shortURL, err := s.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, fmt.Errorf("cannot shorten URL for id [%d]", url.ID)
	}
	log.Printf("Shortened: %s - %s", url, shortURL)

	return shortURL, nil
}
