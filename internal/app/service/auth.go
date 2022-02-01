package service

import (
	"crypto/aes"
	"encoding/hex"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/pkg/idgenerator"
)

type Auth struct {
}

func (auth *Auth) CreateSissionID() (string, string, error) {
	//generate SessionID
	id := idgenerator.CreateID(16)
	src, err := hex.DecodeString(id)
	if err != nil {
		return "", "", err
	}
	//read secret key
	key := []byte(configs.Secret)

	//sign the session whis a secret key
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}
	dst := make([]byte, 16)

	aesblock.Encrypt(dst, src)
	dstStr := hex.EncodeToString(dst)

	return id, dstStr, nil
}

func (auth *Auth) ReadSessionID(id string) (string, error) {
	key := []byte(configs.Secret)
	dst, err := hex.DecodeString(id)
	if err != nil {
		return "", err
	}
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//decryption session id by secret key
	src := make([]byte, 16)
	aesblock.Decrypt(src, dst)
	return hex.EncodeToString(src), nil
}
