package model

import (
	"fmt"
	"net/url"
)

// URLToShorten represents an intent and input data to shorten a long URL
type URLToShorten struct {
	UserID  int64
	LongURL url.URL
}

func NewURLToShorten(userID int64, u url.URL) URLToShorten {
	return URLToShorten{
		UserID:  userID,
		LongURL: u,
	}
}

// ShortenedURL represents a successfully shortened and stored URL
type ShortenedURL struct {
	UserID  int64
	ID      int
	LongURL url.URL
}

func NewShortenedURL(userID int64, id int, u url.URL) ShortenedURL {
	return ShortenedURL{
		UserID:  userID,
		ID:      id,
		LongURL: u,
	}
}

// String provides a text representation of a shortened URL;
// useful for logging
func (u ShortenedURL) String() string {
	return fmt.Sprintf("StoreURL{%d - %d - %s}", u.UserID, u.ID, u.LongURL.String())
}
