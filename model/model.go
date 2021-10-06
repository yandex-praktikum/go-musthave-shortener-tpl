package model

import (
	"fmt"
	"net/url"
)

type StorableURL struct {
	LongURL url.URL
}

func NewStorableURL(u url.URL) StorableURL {
	return StorableURL{u}
}

type StoreURL struct {
	ID      int
	LongURL url.URL
}

func NewStoreURL(id int, u url.URL) StoreURL {
	return StoreURL{id, u}
}

func (u StoreURL) String() string {
	return fmt.Sprintf("StoreURL{%d - %s}", u.ID, u.LongURL.String())
}
