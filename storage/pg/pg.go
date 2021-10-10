package pg

import (
	"database/sql"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type PgStorage struct {
	*sql.DB
}

func New(db *sql.DB) *PgStorage {
	return &PgStorage{db}
}

func (s *PgStorage) GetByID(id int) *model.ShortenedURL {
	return nil
}

func (s *PgStorage) ListByUserID(userID int) []model.ShortenedURL {
	result := make([]model.ShortenedURL, 0)

	return result
}

func (s *PgStorage) Save(model.URLToShorten) model.ShortenedURL {
	return model.ShortenedURL{}
}
