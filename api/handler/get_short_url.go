package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

func (h *URLShortenerHandler) handleGetShortURL(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid id [%v]", idStr)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	url, errGet := h.Service.GetURLByID(id)
	if errors.Is(errGet, model.ErrURLNotFound) {
		http.NotFound(w, r)
		return
	}
	if errGet != nil {
		log.Printf("Cannot get URL by ID [%d]: %s", id, errGet.Error())
		http.Error(w, "Cannot get URL by ID", http.StatusInternalServerError)
		return
	}

	log.Printf("Found: %d - %v", id, url)

	w.Header().Add("Location", url.LongURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
