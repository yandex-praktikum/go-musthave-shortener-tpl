package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

func (h *URLShortenerHandler) handlePostLongURL(w http.ResponseWriter, r *http.Request) {
	rawURL, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("Cannot read request: %v", errRead)
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}

	log.Printf("Got url to shorten: %s", rawURL)
	longURL, errParse := url.Parse(string(rawURL))
	if errParse != nil {
		log.Printf("Cannot parse URL: %v", errParse)
		http.Error(w, "Cannot parse URL", http.StatusBadRequest)
		return
	}

	userID := userID(r)
	u := model.NewURLToShorten(userID, *longURL)
	shortenedURL, errShorten := h.Service.ShortenURL(u)
	if errors.Is(errShorten, model.ErrDuplicateURL) {
		url, errGet := h.Service.LookupURL(u.LongURL)
		if errGet != nil {
			log.Printf("Duplicate URL, but cannot find [%s]: %s", u.LongURL.String(), errGet.Error())
			http.Error(w, "Duplicate URL, but cannot find", http.StatusInternalServerError)
			return
		}

		h.writeResponsePostLongURL(w, http.StatusConflict, *url)
		return
	}
	if errShorten != nil {
		log.Printf("Cannot shorten url: %s", errShorten.Error())
		http.Error(w, "Cannot shorten url", http.StatusInternalServerError)
		return
	}

	h.writeResponsePostLongURL(w, http.StatusCreated, *shortenedURL)
}

func (h *URLShortenerHandler) writeResponsePostLongURL(w http.ResponseWriter, status int, u model.ShortenedURL) {
	absoluteURL, errAbsolute := h.Service.AbsoluteURL(u)
	if errAbsolute != nil {
		log.Printf("Cannot resolve absolute URL: %s", errAbsolute)
		http.Error(w, "Cannot resolve absolute URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, errWrite := fmt.Fprint(w, absoluteURL.String()); errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}
