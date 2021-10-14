// Package apimodel provides structures for (un)marshalling JSON bodies
// of HTTP requests and responses.
package apimodel

type LongURLJson struct {
	URL string `json:"url"`
}

type ShortURLJson struct {
	Result string `json:"result"`
}

type ShortURLForUserJSON struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type LongBatchURLJson struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

type ShortBatchURLJson struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"short_url"`
}
