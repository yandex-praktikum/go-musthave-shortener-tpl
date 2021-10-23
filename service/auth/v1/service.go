package v1

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

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

	userToAdd := model.UserToAdd{Key: hex.EncodeToString(key)}
	user, errAdd := s.SaveUser(userToAdd)
	if errAdd != nil {
		return nil, fmt.Errorf("cannot save new user: %w", errAdd)
	}

	return user, nil
}

func (s *Service) SignUserID(u model.User) (*model.SignedUserID, error) {
	key, errHex := hex.DecodeString(u.Key)
	if errHex != nil {
		return nil, fmt.Errorf("invalid key [%s]: %w", u.Key, errHex)
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(strconv.FormatInt(u.ID, 10)))
	hmac := h.Sum(nil)

	signedUserID := model.SignedUserID{
		ID:   u.ID,
		HMAC: hex.EncodeToString(hmac),
	}

	return &signedUserID, nil
}

func (s *Service) Validate(sgn model.SignedUserID) error {
	u, errGet := s.GetUserByID(sgn.ID)
	if errGet != nil {
		return fmt.Errorf("cannot validate signature: %w", errGet)
	}

	if u.ID != sgn.ID {
		msg := fmt.Sprintf("trying to check signature (ID [%d]) for other user (ID [%d])", sgn.ID, u.ID)
		return errors.New(msg)
	}

	canonicalS, errSign := s.SignUserID(*u)
	if errSign != nil {
		return fmt.Errorf("cannot get signature for user [%d]: %w", u.ID, errSign)
	}

	sgnHMAC, errSgn := hex.DecodeString(sgn.HMAC)
	if errSgn != nil {
		return fmt.Errorf("invalid signed user id HMAC [%s]: %w", sgn.HMAC, errSgn)
	}

	canonicalHMAC, errCanonical := hex.DecodeString(canonicalS.HMAC)
	if errCanonical != nil {
		return fmt.Errorf("invalid canonical user HMAC [%s]: %w", canonicalS.HMAC, errCanonical)
	}

	if !hmac.Equal(canonicalHMAC, sgnHMAC) {
		msg := fmt.Sprintf("signature %v doesn't match", sgn)
		return errors.New(msg)
	}

	return nil
}
