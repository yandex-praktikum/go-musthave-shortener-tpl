package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/apiModel"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

func (h *URLShortenerHandler) handlePostAPIShorten(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	longURLJson := apimodel.LongURLJson{}
	if errDec := dec.Decode(&longURLJson); errDec != nil {
		msg := fmt.Sprintf("Cannot decode request body: %v", errDec)
		http.Error(w, msg, http.StatusBadRequest)
		log.Println(msg)
		return
	}
	longURL, errParse := url.Parse(longURLJson.URL)
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	log.Printf("longURLJson.Url: [%v]", longURL)

	storableURL := model.NewURLToShorten(*longURL)
	shortURL, errShorten := h.Service.ShortenURL(storableURL)
	if errShorten != nil {
		log.Printf("Cannot shorten url: %s", errShorten.Error())
		http.Error(w, "Cannot shorten url", http.StatusInternalServerError)
		return
	}

	shortURLJson := apimodel.ShortURLJson{Result: shortURL.String()}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if errEnc := enc.Encode(shortURLJson); errEnc != nil {
		log.Printf("Cannot write response: %v", errEnc)
	}
}
