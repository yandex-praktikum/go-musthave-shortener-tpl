package mocks

import (
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) GetURLByID(id int) *model.ShortenedURL {
	args := m.Called(id)

	return args.Get(0).(*model.ShortenedURL)
}

func (m *StorageMock) ListByUserID(userID int) []model.ShortenedURL {
	args := m.Called(userID)

	return args.Get(0).([]model.ShortenedURL)
}

func (m *StorageMock) SaveURL(u model.URLToShorten) model.ShortenedURL {
	args := m.Called(u)

	return args.Get(0).(model.ShortenedURL)
}
