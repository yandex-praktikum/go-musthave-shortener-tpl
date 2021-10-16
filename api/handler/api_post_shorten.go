package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	apimodel "github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/apiModel"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
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

	userID := userID(r)
	u := model.NewURLToShorten(userID, *longURL)
	shortenedURL, errShorten := h.Service.ShortenURL(u)
	if errors.Is(errShorten, storage.ErrDuplicateURL) {
		url, errGet := h.Service.LookupURL(u.LongURL)
		if errGet != nil {
			log.Printf("Duplicate URL, but cannot find [%s]: %s", u.LongURL.String(), errGet.Error())
			http.Error(w, "Duplicate URL, but cannot find", http.StatusInternalServerError)
			return
		}

		h.writeResponseAPIPostShorten(w, http.StatusConflict, *url)
		return
	}
	if errShorten != nil {
		log.Printf("Cannot shorten url: %s", errShorten.Error())
		http.Error(w, "Cannot shorten url", http.StatusInternalServerError)
		return
	}

	h.writeResponseAPIPostShorten(w, http.StatusCreated, *shortenedURL)
}

func (h *URLShortenerHandler) writeResponseAPIPostShorten(w http.ResponseWriter, status int, u model.ShortenedURL) {
	absoluteURL, errAbsolute := h.Service.AbsoluteURL(u)
	if errAbsolute != nil {
		log.Printf("Cannot resolve absolute URL: %s", errAbsolute)
		http.Error(w, "Cannot resolve absolute URL", http.StatusInternalServerError)
		return
	}
	shortURLJson := apimodel.ShortURLJson{Result: absoluteURL.String()}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if errEnc := enc.Encode(shortURLJson); errEnc != nil {
		log.Printf("Cannot write response: %v", errEnc)
	}
}
