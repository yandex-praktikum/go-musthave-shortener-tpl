package v1

import (
	"fmt"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type Service struct {
	storage.AuthStorage
}

func New(storage storage.AuthStorage) *Service {
	return &Service{storage}
}

func (s *Service) SignUp() (*model.User, error) {
	key, errKey := auth.GenerateKey()
	if errKey != nil {
		return nil, fmt.Errorf("cannot sign up user: %w", errKey)
	}

	userToAdd := model.UserToAdd{Key: key}
	user, errAdd := s.Save(userToAdd)
	if errAdd != nil {
		return nil, fmt.Errorf("cannot save new user: %w", errAdd)
	}

	return user, nil
}

func (s *Service) Validate(sgn model.SignedUserID) error {
	u, errGet := s.GetByID(sgn.ID)
	if errGet != nil {
		return fmt.Errorf("cannot validate signature: %w", errGet)
	}

	return auth.ValidateSignature(*u, sgn)
}
