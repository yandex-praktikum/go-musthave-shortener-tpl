package model

type URL struct {
	URLID     string `json:"-"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"original_url"`
	SessionID string `json:"-"`
	IsDeleted bool   `json:"-"`
}
