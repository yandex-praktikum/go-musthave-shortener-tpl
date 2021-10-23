package pg

import (
	"database/sql"
	"fmt"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type PgAuthStorage struct {
	*sql.DB
}

func NewAuthStorage(db *sql.DB) *PgAuthStorage {
	return &PgAuthStorage{db}
}

func (s *PgAuthStorage) GetByID(id int64) (*model.User, error) {
	row := s.QueryRow("select USERS_ID, USERS_SIGN_KEY from USERS where USERS_ID = $1", id)
	user := model.User{}

	errScan := row.Scan(&user.ID, &user.Key)
	if errScan != nil {
		return nil, fmt.Errorf("cannot get user by id [%d]: %w", id, errScan)
	}

	return &user, nil
}

func (s *PgAuthStorage) Save(u model.UserToAdd) (*model.User, error) {
	row := s.QueryRow("insert into USERS (USERS_SIGN_KEY) values($1) returning USERS_ID, USERS_SIGN_KEY", u.Key)
	user := model.User{}

	errScan := row.Scan(&user.ID, &user.Key)
	if errScan != nil {
		return nil, fmt.Errorf("cannot insert user: %w", errScan)
	}

	return &user, nil
}
