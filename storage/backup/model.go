package backup

import (
	"fmt"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type gobURL struct {
	ID      int
	LongURL string
}

func newGobURL(u model.StoreURL) gobURL {
	return gobURL{
		ID:      u.ID,
		LongURL: u.LongURL.String(),
	}
}

func (u *gobURL) ToStoreURL() (*model.StoreURL, error) {
	url, errParse := url.Parse(u.LongURL)
	if errParse != nil {
		return nil, fmt.Errorf("cannot restore url [%s] from backup: %w", u.LongURL, errParse)
	}
	storeURL := model.NewStoreURL(u.ID, url)
	return &storeURL, nil
}
