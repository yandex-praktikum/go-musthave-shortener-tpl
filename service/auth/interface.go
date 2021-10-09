package auth

import "github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"

type IDService interface {
	SignUp() (*model.User, error)
	Validate(sgn model.SignedUserID) error
}
