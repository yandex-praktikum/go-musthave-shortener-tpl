package auth

import (
	"crypto/rand"
	"fmt"
)

const KeySize = 16

func GenerateKey() ([]byte, error) {
	b := make([]byte, KeySize)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("cannot generate key: %w", err)
	}

	return b, nil
}
