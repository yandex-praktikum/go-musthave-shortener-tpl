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

	url, errGet := h.Service.GetByID(id)
	if errGet != nil {
		log.Printf("Cannot get URL by ID [%d]: %s", id, errGet.Error())
		http.Error(w, "Cannot get URL by ID", http.StatusInternalServerError)
		return
	}
	if url == nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("Found: %d - %v", id, url)

	w.Header().Add("Location", url.LongURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
