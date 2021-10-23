package mocks

import "github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"

type IDServiceStub struct{}

func NewIDServiceStub() *IDServiceStub {
	return &IDServiceStub{}
}

func (m *IDServiceStub) SignUp() (*model.User, error) {
	return &model.User{ID: 0, Key: make([]byte, 0)}, nil
}

func (m *IDServiceStub) Validate(sgn model.SignedUserID) error {
	return nil
}

func (m *IDServiceStub) SignUserID(u model.User) (*model.SignedUserID, error) {
	return &model.SignedUserID{
		ID:        0,
		Signature: "bb0cb0e08a30cbcbbb4d1cc8ce2bed4dff036f49894ef8b1e0eba1909368ee4b",
	}, nil
}
