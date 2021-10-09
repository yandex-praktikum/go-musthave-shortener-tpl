package mocks

import (
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) GetByID(id int) *model.ShortenedURL {
	args := m.Called(id)

	return args.Get(0).(*model.ShortenedURL)
}

func (m *StorageMock) Save(u model.URLToShorten) model.ShortenedURL {
	args := m.Called(u)

	return args.Get(0).(model.ShortenedURL)
}
