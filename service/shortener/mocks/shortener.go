package mocks

import (
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/stretchr/testify/mock"
)

type URLServiceMock struct {
	mock.Mock
}

func (m *URLServiceMock) ShortenURL(u model.URLToShorten) (*model.ShortenedURL, error) {
	args := m.Called(u)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.ShortenedURL), args.Error(1)
}

func (m *URLServiceMock) GetByID(id int) (*model.ShortenedURL, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.ShortenedURL), args.Error(1)
}

func (m *URLServiceMock) LookupURL(u url.URL) (*model.ShortenedURL, error) {
	args := m.Called(u)

	return args.Get(0).(*model.ShortenedURL), args.Error(1)
}

func (m *URLServiceMock) GetUserURLs(userID int64) ([]model.ShortenedURL, error) {
	args := m.Called(userID)

	return args.Get(0).([]model.ShortenedURL), args.Error(1)
}

func (m *URLServiceMock) AbsoluteURL(u model.ShortenedURL) (*url.URL, error) {
	args := m.Called(u)

	return args.Get(0).(*url.URL), args.Error(1)
}
