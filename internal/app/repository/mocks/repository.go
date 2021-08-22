package mocks

import (
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) SaveShortener(code string, shortener *model.Shortener) error {
	return nil
}

func (m *RepositoryMock) GetShortenerBy(id string) (*model.Shortener, error) {
	shortener := &model.Shortener{
		URL: "https://yandex.ru/",
	}
	return shortener, nil
}
