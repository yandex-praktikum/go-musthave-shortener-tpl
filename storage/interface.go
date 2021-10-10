// Package storage provides a persistent storage for the service
package storage

import "github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"

// ShortenerStorage provides methods to persist and retrieve shortened URLs
// of a single request.
type ShortenerStorage interface {
	// GetByID looks-up for a previously shortened URL.
	GetByID(id int) *model.ShortenedURL

	// ListByUserID returns all URLs shortened by the specified user.
	ListByUserID(userID int) []model.ShortenedURL

	// Save persists a shortened URL. It is responsible for generating ID.
	Save(model.URLToShorten) model.ShortenedURL
}

// AuthStorage provides methods to persist and retrieve
// authentication-related staff.
type AuthStorage interface {
	// GetByID looks-up an existing user
	GetByID(id int) (*model.User, error)

	// Save adds a new user. This method is responsible for generation
	// of a user ID.
	Save(u model.UserToAdd) (*model.User, error)
}
