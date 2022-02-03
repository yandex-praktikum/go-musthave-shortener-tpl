package model

type BatchRequest struct {
	CorrelationID string `json:"correlation_id" binding:"required"`
	OriginalURL   string `json:"original_url" binding:"required"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id" binding:"required"`
	ShortURL      string `json:"short_url" binding:"required"`
}
