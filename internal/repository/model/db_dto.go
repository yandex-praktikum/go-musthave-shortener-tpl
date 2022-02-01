package model

type ShortenDTO struct {
	URLID    string `json:"-"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}
