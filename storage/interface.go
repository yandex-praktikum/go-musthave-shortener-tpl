// Package storage provides a persistent storage for the service
package storage

import "github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"

// Storage provides methods to perform actions in context
// of a single request.
type Storage interface {
	// GetByID looks-up for a previously shortened URL.
	GetByID(id int) *model.ShortenedURL

	// ListByUserID returns all URLs shortened by the specified user.
	ListByUserID(userID int) []model.ShortenedURL

	// Save persists a shortened URL. It is responsible for generating ID.
	Save(model.URLToShorten) model.ShortenedURL
}

// BulkStorage is used to backup and restore the server state
// It extends Storage for convenience because usually it is the same
// instance.
type BulkStorage interface {
	Storage

	// GetAll returns all shortened URLs.
	GetAll() []model.ShortenedURL

	// Load persists a shortened URL. It is different from Storage.Save
	// in that is does not generate an ID, but takes whatever comes
	// with the URL.
	Load(u model.ShortenedURL)
}
