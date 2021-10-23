package handler

import (
	"log"
	"net/http"
)

func (h *URLShortenerHandler) handleGetPing(w http.ResponseWriter, r *http.Request) {
	if err := h.Pinger.Ping(); err != nil {
		log.Printf("Ping failed: %s", err.Error())
		http.Error(w, "Ping failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
