package v1

import (
	"errors"
	"fmt"
	"sync"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type Service struct {
	sync.RWMutex

	users []model.User
}

func New(_storage storage.Storage) *Service {
	return &Service{
		RWMutex: sync.RWMutex{},
		users:   make([]model.User, 0),
	}
}

func (s *Service) SignUp() (*model.User, error) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	key, errKey := auth.GenerateKey()
	if errKey != nil {
		return nil, fmt.Errorf("cannot sign up user: %w", errKey)
	}

	id := len(s.users)
	u := model.User{
		ID:  id,
		Key: key,
	}
	s.users = append(s.users, u)

	return &u, nil
}

func (s *Service) Validate(sgn model.SignedUserID) error {
	for _, u := range s.users {
		if u.ID == sgn.ID {
			return auth.ValidateSignature(u, sgn)
		}
	}

	msg := fmt.Sprintf("cannot find user with ID [%d]", sgn.ID)

	return errors.New(msg)
}
