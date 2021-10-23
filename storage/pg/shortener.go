package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type PgShortenerStorage struct {
	*sql.DB
}

func NewShortenerStorage(db *sql.DB) (*PgShortenerStorage, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &PgShortenerStorage{db}, nil
}

func (s *PgShortenerStorage) GetURLByID(id int) (*model.ShortenedURL, error) {
	row := s.QueryRow("select URLS_ID, URLS_ORIGINAL_URL, USERS_ID from URLS where URLS_ID = $1", id)

	url := model.ShortenedURL{}

	if err := mapShortenedURL(&url, row); err != nil {
		return nil, fmt.Errorf("cannot get URL by id [%d]: %w", id, err)
	}

	return &url, nil
}

func (s *PgShortenerStorage) LookupURL(u url.URL) (*model.ShortenedURL, error) {
	row := s.QueryRow("select URLS_ID, URLS_ORIGINAL_URL, USERS_ID from URLS where URLS_ORIGINAL_URL = $1", u.String())

	url := model.ShortenedURL{}

	if err := mapShortenedURL(&url, row); err != nil {
		return nil, fmt.Errorf("cannot lookup URL [%s]: %w", u.String(), err)
	}

	return &url, nil
}

func (s *PgShortenerStorage) ListByUserID(userID int64) ([]model.ShortenedURL, error) {
	result := make([]model.ShortenedURL, 0)

	rows, err := s.Query(`
		select URLS_ID, URLS_ORIGINAL_URL, USERS_ID
		from URLS
		where USERS_ID = $1
	`,
		userID)
	if err != nil {
		return result, fmt.Errorf("cannot select URLs for user [%d]: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		url := model.ShortenedURL{}

		if err := mapShortenedURL(&url, rows); err != nil {
			return result, fmt.Errorf("cannot map all urls from DB: %w", err)
		}

		result = append(result, url)
	}
	if rows.Err() != nil {
		return result, fmt.Errorf("cannot iterate all results from DB: %w", rows.Err())
	}

	return result, nil
}

func (s *PgShortenerStorage) SaveURL(u model.URLToShorten) (model.ShortenedURL, error) {
	row := s.QueryRow(`
		insert into URLS (URLS_ORIGINAL_URL, USERS_ID) 
		values($1, $2)
		returning URLS_ID, URLS_ORIGINAL_URL, USERS_ID
	`, u.LongURL.String(), u.UserID)

	url := model.ShortenedURL{}
	if err := mapShortenedURL(&url, row); err != nil {
		var dbErr *pgconn.PgError
		if errors.As(err, &dbErr) && dbErr.Code == pgerrcode.UniqueViolation {
			log.Printf("Duplicate URL: %s", u.LongURL.String())
			err = model.ErrDuplicateURL
		}
		return url, fmt.Errorf("cannot insert url: %w", err)
	}

	return url, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func mapShortenedURL(u *model.ShortenedURL, row scannable) error {
	var longURLStr string

	errScan := row.Scan(&u.ID, &longURLStr, &u.UserID)
	if errScan == sql.ErrNoRows {
		return model.ErrURLNotFound
	}
	if errScan != nil {
		return fmt.Errorf("cannot scan url from DB results: %w", errScan)
	}

	longURL, err := url.Parse(longURLStr)
	if err != nil {
		return fmt.Errorf("invalid URL [%s]: %w", longURLStr, err)
	}

	u.LongURL = *longURL

	return nil
}
