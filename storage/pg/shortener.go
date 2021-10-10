package pg

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type PgShortenerStorage struct {
	*sql.DB
}

func NewShortenerStorage(db *sql.DB) *PgShortenerStorage {
	return &PgShortenerStorage{db}
}

func (s *PgShortenerStorage) GetByID(id int) (*model.ShortenedURL, error) {
	row := s.QueryRow("select URLS_ID, URLS_ORIGINAL_URL, USERS_ID from URLS where URLS_ID = $1", id)

	url, errMap := mapShortenedURL(row)
	if errMap != nil {
		return nil, fmt.Errorf("cannot get url by id [%d]: %w", id, errMap)
	}

	return url, nil
}

func (s *PgShortenerStorage) ListByUserID(userID int) ([]model.ShortenedURL, error) {
	result := make([]model.ShortenedURL, 0)

	rows, errQuery := s.Query(`
		select URLS_ID, URLS_ORIGINAL_URL, USERS_ID
		from URLS
		where USERS_ID = $1
	`,
		userID)
	if errQuery != nil {
		return result, fmt.Errorf("cannot select URLs for user [%d]: %w", userID, errQuery)
	}
	defer rows.Close()

	for rows.Next() {
		url, errMap := mapShortenedURL(rows)
		if errMap != nil {
			return result, fmt.Errorf("cannot map all urls from DB: %w", errMap)
		}

		result = append(result, *url)
	}
	if rows.Err() != nil {
		return result, fmt.Errorf("cannot iterate all results from DB: %w", rows.Err())
	}

	return result, nil
}

func (s *PgShortenerStorage) Save(u model.URLToShorten) (*model.ShortenedURL, error) {
	row := s.QueryRow(`
		insert into URLS (URLS_ORIGINAL_URL, USERS_ID) 
		values($1, $2) 
		returning URLS_ID, URLS_ORIGINAL_URL, USERS_ID
	`, u.LongURL.String(), u.UserID)

	url, errMap := mapShortenedURL(row)
	if errMap != nil {
		return nil, fmt.Errorf("cannot insert url: %w", errMap)
	}

	return url, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func mapShortenedURL(row scannable) (*model.ShortenedURL, error) {
	var id, userID int
	var longURLStr string

	errScan := row.Scan(&id, &longURLStr, &userID)
	if errScan == sql.ErrNoRows {
		return nil, nil
	}
	if errScan != nil {
		return nil, fmt.Errorf("cannot scan url from DB results: %w", errScan)
	}

	longURL, errParse := url.Parse(longURLStr)
	if errParse != nil {
		return nil, fmt.Errorf("invalid URL [%s]: %w", longURLStr, errParse)
	}

	shortenedURL := model.ShortenedURL{
		ID:      id,
		UserID:  userID,
		LongURL: *longURL,
	}

	return &shortenedURL, nil
}
