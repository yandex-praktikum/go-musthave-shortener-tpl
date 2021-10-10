package pg

import (
	"database/sql"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type PgAuthStorage struct {
	*sql.DB
}

func NewAuthStorage(db *sql.DB) *PgAuthStorage {
	return &PgAuthStorage{db}
}

func (s *PgAuthStorage) GetByID(id int) (*model.User, error) {
	return nil, nil
}

func (s *PgAuthStorage) Save(model.UserToAdd) (*model.User, error) {
	return nil, nil
}
