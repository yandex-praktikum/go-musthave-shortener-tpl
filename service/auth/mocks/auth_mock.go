package mocks

import (
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/stretchr/testify/mock"
)

type IDServiceMock struct {
	mock.Mock
}

func (m *IDServiceMock) SignUp() (*model.User, error) {
	args := m.Called()

	return args.Get(0).(*model.User), args.Error(1)
}

func (m *IDServiceMock) Validate(sgn model.SignedUserID) error {
	args := m.Called(sgn)

	return args.Error(1)
}

func (m *IDServiceMock) SignUserID(u model.User) (*model.SignedUserID, error) {
	args := m.Called()

	return args.Get(0).(*model.SignedUserID), args.Error(1)
}
