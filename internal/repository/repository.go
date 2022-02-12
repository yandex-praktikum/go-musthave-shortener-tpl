package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Storage struct {
	client Client
	*DeleteBuffer
	ShutdownChan    chan struct{}
	DeleteCompleted chan struct{}
}

func NewStorage(client Client) *Storage {
	return &Storage{
		client:          client,
		DeleteBuffer:    NewDeleteBuffer(),
		ShutdownChan:    make(chan struct{}),
		DeleteCompleted: make(chan struct{}),
	}
}

//=====================================================================================================
func (us *Storage) SaveURL(m *model.URL, key string) (string, error) {
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

//=====================================================================================================
func (us *Storage) SaveBatch(list *[]model.URL, key string) error {
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

//=====================================================================================================
func (us *Storage) GetURL(key string) (string, error) {
	var longURL string
	var isDeleted bool
	q := `SELECT long_url,is_deleted FROM shortens
	WHERE
		url_id=$1;`
	us.client.QueryRow(context.Background(), q, key).Scan(&longURL, &isDeleted)
	if isDeleted {
		err := errors.New("error: The URL has been deleted")
		return "", fmt.Errorf(`%w`, err)
	}
	if longURL == "" {
		return "", errors.New("error: Not found data")
	}

	return longURL, nil
}

//=====================================================================================================
func (us *Storage) PingDB() error {
	return us.client.Ping(context.Background())
}

//=====================================================================================================
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

//=====================================================================================================
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

//=====================================================================================================
func (us *Storage) GetList(key string) ([]model.URL, error) {
	var list []model.URL
	q := `SELECT short_url, long_url 
		FROM shortens 
	WHERE 
		session_id=(select id from sessions where session_id =$1)`
	rows, err := us.client.Query(context.Background(), q, key)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var model model.URL
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

//=====================================================================================================
func (us *Storage) DeleteURLs(list []model.URL) {
	if len(us.DeleteBuffer.Buffer) == 0 {
		go func() {
			us.DeleteCompleted <- struct{}{}
		}()
		return
	}

	q := `UPDATE shortens
		SET is_deleted=true 
		WHERE
		url_id=$1 AND
		session_id=(select id from sessions where session_id =$2);`

	batch := &pgx.Batch{}
	for _, val := range list {
		batch.Queue(q, val.URLID, val.SessionID)
	}
	br := us.client.SendBatch(context.Background(), batch)
	defer br.Close()
	br.Exec()
	go func() {
		us.DeleteCompleted <- struct{}{}
	}()
}

//=====================================================================================================
func (us *Storage) AddToBuffer(m model.URL) {
	us.DeleteBuffer.Mutex.Lock()

	//add item to buffer
	us.DeleteBuffer.Buffer = append(us.DeleteBuffer.Buffer, m)

	go func() {
		time.Sleep(time.Second * 10)
		go us.SentMessage("timer")
		us.DeleteBuffer.Signal <- struct{}{}
	}()

	if cap(us.DeleteBuffer.Buffer) == len(us.DeleteBuffer.Buffer) {
		go us.SentMessage("overflow")
		us.DeleteBuffer.Signal <- struct{}{}
	}

	us.DeleteBuffer.Mutex.Unlock()
}

//=====================================================================================================
func (us *Storage) DeleteBufferRefreshing() {
	for {
		select {
		case <-us.DeleteBuffer.Signal:
			//if buffer is full -> sent to db
			msg := <-us.DeleteBuffer.Msg
			us.DeleteURLs(us.DeleteBuffer.Buffer)
			<-us.DeleteCompleted
			//clear buffer
			us.DeleteBuffer.ClearBuffer()
			log.Println("Buffer was cleared by", msg)

		case <-us.ShutdownChan:
			us.DeleteURLs(us.DeleteBuffer.Buffer)
			us.DeleteBuffer.ClearBuffer()
			log.Println("Buffer was cleared by shutdown")
		}
	}
}

//=====================================================================================================
func (us *Storage) SentMessage(msg string) {
	us.DeleteBuffer.Msg <- msg
}
