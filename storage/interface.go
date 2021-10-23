// Package storage provides a persistent storage for the service
package storage

import (
	"errors"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

// ShortenerStorage provides methods to persist and retrieve shortened URLs
// of a single request.
type ShortenerStorage interface {
	// GetURLByID looks-up for a previously shortened URL.
	GetURLByID(id int) (*model.ShortenedURL, error)

	// LookupURL does a search of a shortened URL by long URL.
	LookupURL(u url.URL) (*model.ShortenedURL, error)

	// ListByUserID returns all URLs shortened by the specified user.
	ListByUserID(userID int64) ([]model.ShortenedURL, error)

	// SaveURL persists a shortened URL. It is responsible for generating ID.
	SaveURL(model.URLToShorten) (*model.ShortenedURL, error)
}

var ErrDuplicateURL = errors.New("URL already shortened")

// AuthStorage provides methods to persist and retrieve
// authentication-related staff.
type AuthStorage interface {
	// GetUserByID looks-up an existing user
	GetUserByID(id int64) (*model.User, error)

	// SaveUser adds a new user. This method is responsible for generation
	// of a user ID.
	SaveUser(u model.UserToAdd) (*model.User, error)
}
