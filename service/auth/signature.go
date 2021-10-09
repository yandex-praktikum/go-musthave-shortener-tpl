package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

func SignUserID(u model.User) model.SignedUserID {
	h := hmac.New(sha256.New, u.Key)
	h.Write([]byte(strconv.Itoa(u.ID)))
	hmac := h.Sum(nil)

	return model.SignedUserID{
		ID:   u.ID,
		HMAC: hmac,
	}
}

func ValidateSignature(u model.User, sgn model.SignedUserID) error {
	if u.ID != sgn.ID {
		msg := fmt.Sprintf("trying to check signature (ID [%d]) for other user (ID [%d])", sgn.ID, u.ID)
		return errors.New(msg)
	}

	canonicalS := SignUserID(u)
	if !hmac.Equal(canonicalS.HMAC, sgn.HMAC) {
		msg := fmt.Sprintf("signature %v doesn't match", sgn)
		return errors.New(msg)
	}

	return nil
}
