package handler

import (
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

	u := model.NewURLToShorten(*longURL)
	shortenedURL, errShorten := h.Service.ShortenURL(u)
	if errShorten != nil {
		log.Printf("Cannot shorten url: %s", errShorten.Error())
		http.Error(w, "Cannot shorten url", http.StatusInternalServerError)
		return
	}

	absoluteURL, errAbsolute := h.Service.AbsoluteURL(*shortenedURL)
	if errAbsolute != nil {
		log.Printf("Cannot resolve absolute URL: %s", errAbsolute)
		http.Error(w, "Cannot resolve absolute URL", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	if _, errWrite := fmt.Fprint(w, absoluteURL.String()); errWrite != nil {
		log.Printf("Cannot write response: %v", errWrite)
	}
}
