package shortener

import (
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type URLShortener interface {
	ShortenURL(newURL model.URLToShorten) (*url.URL, error)
}
