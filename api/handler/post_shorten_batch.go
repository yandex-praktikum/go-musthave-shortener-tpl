package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	apimodel "github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/apiModel"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

func (h *URLShortenerHandler) handlePostShortenBatch(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	longURLs := make([]apimodel.LongBatchURLJson, 0)

	if errDecode := dec.Decode(&longURLs); errDecode != nil {
		log.Printf("Cannot read request: %s", errDecode.Error())
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}

	userID := userID(r)
	shortURLs := make([]apimodel.ShortBatchURLJson, 0)
	for _, longBatchURL := range longURLs {
		longURL, errParse := url.Parse(longBatchURL.URL)
		if errParse != nil {
			msg := fmt.Sprintf("Bad url with correlation id [%s]", longBatchURL.CorrelationID)
			log.Printf("%s: %v", msg, errParse)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		u := model.NewURLToShorten(userID, *longURL)
		shortURL, errShorten := h.Service.ShortenURL(u)
		if errShorten != nil {
			msg := fmt.Sprintf("Cannot shorten url with correlation id [%s]", longBatchURL.CorrelationID)
			log.Printf("%s: %v", msg, errShorten.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		absURL, errAbs := h.Service.AbsoluteURL(*shortURL)
		if errAbs != nil {
			msg := fmt.Sprintf("Cannot make absolute url;correlation id [%s]", longBatchURL.CorrelationID)
			log.Printf("%s: %v", msg, errAbs.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		shortBatchURL := apimodel.ShortBatchURLJson{
			CorrelationID: longBatchURL.CorrelationID,
			URL:           absURL.String(),
		}

		shortURLs = append(shortURLs, shortBatchURL)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if errEnc := enc.Encode(shortURLs); errEnc != nil {
		log.Printf("Cannot write response: %s", errEnc.Error())
		return
	}
}
