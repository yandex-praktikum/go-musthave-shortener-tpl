package v1

import (
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

func New(s storage.ShortenerStorage, u url.URL) *Service {
	return &Service{s, u}
}

func (s *Service) ShortenURL(u model.URLToShorten) (*model.ShortenedURL, error) {
	url := s.Storage.Save(u)
	log.Printf("Shortened: %s", url)

	return &url, nil
}

func (s *Service) AbsoluteURL(u model.ShortenedURL) (*url.URL, error) {
	urlPath := fmt.Sprintf("%d", u.ID)

	shortURL, err := s.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, fmt.Errorf("cannot make absolute URL for id [%d]", u.ID)
	}

	return shortURL, nil
}
