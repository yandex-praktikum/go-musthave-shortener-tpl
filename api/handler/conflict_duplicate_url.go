package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

func (h *URLShortenerHandler) handleConflictDuplicateUrl(w http.ResponseWriter, u model.URLToShorten) {
	url, errGet := h.Service.LookupURL(u.LongURL)
	if errGet != nil {
		log.Printf("Duplicate URL, but cannot find [%s]: %s", u.LongURL.String(), errGet.Error())
		http.Error(w, "Duplicate URL, but cannot find", http.StatusInternalServerError)
		return
	}
	log.Printf("Found: %d - %v", url.ID, url)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	enc := json.NewEncoder(w)
	if errEnc := enc.Encode(url); errEnc != nil {
		log.Printf("Cannot write response: %v", errEnc)
	}
}
