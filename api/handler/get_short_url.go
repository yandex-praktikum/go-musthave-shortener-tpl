package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *URLShortenerHandler) handleGetShortURL(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid id [%v]", idStr)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	url := h.Storage.GetByID(id)
	if url == nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("Found: %d - %v", id, url)

	w.Header().Add("Location", url.LongURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
