package handler

import (
	"net/http"
)

func (h *URLShortenerHandler) handleGetPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
