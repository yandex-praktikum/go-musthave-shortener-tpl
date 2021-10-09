package model

import (
	"fmt"
	"net/url"
)

// URLToShorten represents an intent and input data to shorten a long URL
type URLToShorten struct {
	LongURL url.URL
}

func NewURLToShorten(u url.URL) URLToShorten {
	return URLToShorten{u}
}

// ShortenedURL represents a successfully shortened and stored URL
type ShortenedURL struct {
	ID      int
	LongURL url.URL
}

func NewShortenedURL(id int, u url.URL) ShortenedURL {
	return ShortenedURL{id, u}
}

// String provides a text representation of a shortened URL;
// useful for logging
func (u ShortenedURL) String() string {
	return fmt.Sprintf("StoreURL{%d - %s}", u.ID, u.LongURL.String())
}
