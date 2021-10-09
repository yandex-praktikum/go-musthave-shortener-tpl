package auth

import (
	"crypto/rand"
	"fmt"
)

const KEY_SIZE = 16

func GenerateKey() ([]byte, error) {
	b := make([]byte, KEY_SIZE)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("cannot generate key: %w", err)
	}

	return b, nil
}
