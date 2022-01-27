package repository

import (
	"context"
	"fmt"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	DBCon pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{DBCon: *pool}
}

func (us *Storage) SaveURL(m *model.Shorten) (string, error) {
	var shortURL string
	q := `INSERT INTO shortens
			(url_id,short_url,long_url,session_id)
	VALUES
			($1,$2,$3,$4)
	RETURNING short_url;`
	us.DBCon.QueryRow(context.Background(), q, m.URLID, m.ShortURL, m.LongURL, m.SessionID).Scan(&shortURL)
	return shortURL, nil
}

func (us *Storage) GetURL(key string) string {
	var longURL string
	q := `SELECT long_url FROM shortens
	WHERE
		url_id=$1;`
	us.DBCon.QueryRow(context.Background(), q, key).Scan(&longURL)
	fmt.Println(key)
	fmt.Println(longURL)
	return longURL
}

func (us *Storage) PingDB() error {
	return us.DBCon.Ping(context.Background())
}

func (us *Storage) SaveCookie(s string) (int, error) {
	var id int
	q := `INSERT INTO sessions(session_id)
    VALUES($1)
	RETURNING id;`
	us.DBCon.QueryRow(context.Background(), q, s).Scan(&id)
	return id, nil
}

func (us *Storage) GetCookieID(s string) int {
	var id int
	q := `SELECT id FROM sessions
	 WHERE 
	session_id= $1;`
	us.DBCon.QueryRow(context.Background(), q, s).Scan(id)
	return id
}
