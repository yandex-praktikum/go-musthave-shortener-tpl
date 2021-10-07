package model

import (
	"fmt"
	"net/url"
)

type URLToShorten struct {
	LongURL url.URL
}

func NewURLToShorten(u url.URL) URLToShorten {
	return URLToShorten{u}
}

type ShortenedURL struct {
	ID      int
	LongURL url.URL
}

func NewShortenedURL(id int, u url.URL) ShortenedURL {
	return ShortenedURL{id, u}
}

func (u ShortenedURL) String() string {
	return fmt.Sprintf("StoreURL{%d - %s}", u.ID, u.LongURL.String())
}
