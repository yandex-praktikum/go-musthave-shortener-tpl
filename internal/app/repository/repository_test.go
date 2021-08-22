package repository

import (
	"errors"
	"testing"

	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/repository/mocks"
)

func TestGetShortenerBy(t *testing.T) {
	repoMock := new(mocks.RepositoryMock)

	repoMock.On("SaveShortener", "testtest", &model.Shortener{URL: "new"}).Return(nil)
	repoMock.AssertExpectations(t)
	repoMock.On("SaveShortener", "", &model.Shortener{URL: "new"}).Return(nil)

	repoMock.On("GetShortenerBy", "testtest").Return("name")
	repoMock.On("GetShortenerBy", "").Return("", errors.New("not found"))
}
