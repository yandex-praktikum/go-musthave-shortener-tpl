package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Storage struct {
	client Client
}

func NewStorage(client Client) *Storage {
	return &Storage{client: client}
}

func (us *Storage) SaveURL(m *model.Shorten, key string) (string, error) {
	var id int
	var shortURL string
	q := `INSERT INTO shortens
	  (url_id,short_url,long_url,session_id)
	  VALUES
	  ($1,$2,$3,
	  (SELECT id FROM sessions WHERE session_id=$4))
	  ON CONFLICT (long_url) 
	  DO UPDATE SET 
	  long_url=EXCLUDED.long_url
	  RETURNING id,short_url;`

	us.client.QueryRow(context.Background(), q, m.URLID, m.ShortURL, m.LongURL, key).Scan(&id, &shortURL)
	if id == 0 {
		return "", errors.New("Internal error: Data was not saved")
	}
	if shortURL != m.ShortURL {
		err := errors.New("Error: Attemt to save data, data already exist")
		log.Println(err)
		return shortURL, fmt.Errorf(`%w`, err)
	}
	return shortURL, nil
}
func (us *Storage) SaveBatch(list *[]model.Shorten, key string) error {
	q := `INSERT INTO shortens
	  (url_id,short_url,long_url,session_id)
	  VALUES
	  ($1,$2,$3,
	  (SELECT id FROM sessions WHERE session_id=$4));`

	batch := &pgx.Batch{}
	for _, val := range *list {
		batch.Queue(q, val.URLID, val.ShortURL, val.LongURL, key)
	}
	br := us.client.SendBatch(context.Background(), batch)
	_, err := br.Exec()
	if err != nil {
		return err
	}
	br.Close()

	return nil
}

func (us *Storage) GetURL(key string) string {
	var longURL string
	q := `SELECT long_url FROM shortens
	WHERE
		url_id=$1;`
	us.client.QueryRow(context.Background(), q, key).Scan(&longURL)
	return longURL
}

func (us *Storage) PingDB() error {
	return us.client.Ping(context.Background())
}

func (us *Storage) SaveCookie(s string) error {
	var id int
	q := `INSERT INTO sessions(session_id)
    VALUES($1)
	RETURNING id;`
	us.client.QueryRow(context.Background(), q, s).Scan(&id)
	if id == 0 {
		return errors.New("Cookie is unsaved")
	}
	return nil
}

func (us *Storage) GetCookie(s string) error {
	var id int
	q := `SELECT id FROM sessions
	 	WHERE 
	session_id= $1;`
	us.client.QueryRow(context.Background(), q, s).Scan(&id)
	if id == 0 {
		return errors.New("Error: Cookie not found")
	}
	return nil
}

func (us *Storage) GetList(key string) ([]model.Shorten, error) {
	var list []model.Shorten
	q := `SELECT short_url, long_url 
		FROM shortens 
	WHERE 
		session_id=(select id from sessions where session_id =$1)`
	rows, err := us.client.Query(context.Background(), q, key)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var model model.Shorten
		err := rows.Scan(&model.ShortURL, &model.LongURL)
		if err != nil {
			return nil, err
		}
		list = append(list, model)
	}
	if len(list) == 0 {
		return nil, errors.New("Not foud data")
	}
	return list, nil
}
