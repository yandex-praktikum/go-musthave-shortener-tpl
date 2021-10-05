package backup

import (
	"fmt"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type GobURL struct {
	ID      int
	LongURL string
}

func NewGobUrl(u model.StoreURL) GobURL {
	return GobURL{
		ID:      u.ID,
		LongURL: u.LongURL.String(),
	}
}

func (u *GobURL) ToStoreURL() (*model.StoreURL, error) {
	url, errParse := url.Parse(u.LongURL)
	if errParse != nil {
		return nil, fmt.Errorf("cannot restore url [%s] from backup: %w", u.LongURL, errParse)
	}
	storeUrl := model.NewStoreURL(u.ID, url)
	return &storeUrl, nil
}
