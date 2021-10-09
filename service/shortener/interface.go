package shortener

import (
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type URLShortener interface {
	// ShortenURL does all the necessary logic and stores produced short URL
	ShortenURL(u model.URLToShorten) (*model.ShortenedURL, error)

	// AbsoluteURL resolves a short URL with regards to base URL
	AbsoluteURL(u model.ShortenedURL) (*url.URL, error)
}
