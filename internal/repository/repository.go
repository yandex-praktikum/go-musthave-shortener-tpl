package repository

import (
	"context"
	"log"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Storage struct {
	DBCon *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{DBCon: pool}
}

func (us *Storage) SaveURL(m *model.ShortenDTO, key string) error {
	var id int
	q := `insert into shortens
	  (url_id,short_url,long_url,session_id)
	  values
	  ($1,$2,$3,
	  (select id from sessions where session_id=$4))
	  RETURNING id;`

	us.DBCon.QueryRow(context.Background(), q, m.URLID, m.ShortURL, m.LongURL, key).Scan(&id)
	if id == 0 {
		return errors.New("Data was not saved")
	}
	return nil
}
func (us *Storage) SaveBatch(list *[]model.ShortenDTO, key string) error {
	log.Println(key)
	q := `insert into shortens
	  (url_id,short_url,long_url,session_id)
	  values
	  ($1,$2,$3,
	  (select id from sessions where session_id=$4));`

	batch := &pgx.Batch{}
	for _, val := range *list {
		batch.Queue(q, val.URLID, val.ShortURL, val.LongURL, key)
	}
	br := us.DBCon.SendBatch(context.Background(), batch)
	_, err := br.Exec()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (us *Storage) GetURL(key string) string {
	var longURL string
	q := `SELECT long_url FROM shortens
	WHERE
		url_id=$1;`
	us.DBCon.QueryRow(context.Background(), q, key).Scan(&longURL)
	return longURL
}

func (us *Storage) PingDB() error {
	return us.DBCon.Ping(context.Background())
}

func (us *Storage) SaveCookie(s string) error {
	var id int
	q := `INSERT INTO sessions(session_id)
    VALUES($1)
	RETURNING id;`
	us.DBCon.QueryRow(context.Background(), q, s).Scan(&id)
	if id == 0 {
		return errors.New("Cookie is unsaved")
	}
	return nil
}

func (us *Storage) GetCookie(s string) bool {
	var isFound bool
	var id int
	q := `SELECT id FROM sessions
	 	WHERE 
	session_id= $1;`
	us.DBCon.QueryRow(context.Background(), q, s).Scan(&id)
	if id != 0 {
		isFound = true
	}
	return isFound
}

func (us *Storage) GetList(key string) ([]model.ShortenDTO, error) {
	var list []model.ShortenDTO
	q := `SELECT short_url, long_url 
		FROM shortens 
	WHERE 
		session_id=(select id from sessions where session_id =$1)`
	rows, err := us.DBCon.Query(context.Background(), q, key)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var model model.ShortenDTO

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
