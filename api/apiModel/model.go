// Package apimodel provides structures for (un)marshalling JSON bodies
// of HTTP requests and responses.
package apimodel

type LongURLJson struct {
	URL string `json:"url"`
}

type ShortURLJson struct {
	Result string `json:"result"`
}
