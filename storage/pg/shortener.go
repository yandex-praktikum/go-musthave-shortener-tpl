package pg

import (
	"database/sql"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type PgShortenerStorage struct {
	*sql.DB
}

func NewShortenerStorage(db *sql.DB) *PgShortenerStorage {
	return &PgShortenerStorage{db}
}

func (s *PgShortenerStorage) GetByID(id int) *model.ShortenedURL {
	return nil
}

func (s *PgShortenerStorage) ListByUserID(userID int) []model.ShortenedURL {
	result := make([]model.ShortenedURL, 0)

	return result
}

func (s *PgShortenerStorage) Save(model.URLToShorten) model.ShortenedURL {
	return model.ShortenedURL{}
}
