package helper

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

const defaultCodeLength = 8

type GeneratedString string

// NewStringGenerator creates a new StringGenerator.
func NewGeneratedString() (GeneratedString, error) {
	value, err := generateRandomString(defaultCodeLength)
	if err != nil {
		return "", err
	}
	return GeneratedString(value), nil
}

func generateRandomString(length int) (string, error) {
	if length < 1 {
		return "", errors.New("invalid code length provided")
	}
	rand.Seed(time.Now().UnixNano())
	chars := []rune(
		"abcdefghijkmnpqrstuvwxyz" +
			"123456789",
	)
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String(), nil
}
