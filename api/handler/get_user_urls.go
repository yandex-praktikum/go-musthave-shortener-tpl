package handler

import (
	"encoding/json"
	"log"
	"net/http"

	apimodel "github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/apiModel"
)

func (h *URLShortenerHandler) handleGetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID := userID(r)
	urls, errGet := h.Service.GetUserURLs(userID)
	if errGet != nil {
		log.Printf("Cannot get URLs for user [%d]: %s", userID, errGet.Error())
		http.Error(w, "Cannot list user URLs", http.StatusInternalServerError)
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]apimodel.ShortURLForUserJSON, 0, len(urls))
	for _, url := range urls {
		shortURL, errAbs := h.Service.AbsoluteURL(url)
		if errAbs != nil {
			log.Printf("Cannot render absolute URL for shortened URL [%d]: %s", url.ID, errAbs.Error())
			continue
		}
		response = append(response, apimodel.ShortURLForUserJSON{
			ShortURL:    shortURL.String(),
			OriginalURL: url.LongURL.String(),
		})
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if errEnc := enc.Encode(response); errEnc != nil {
		log.Printf("Cannot write response: %v", errEnc)
	}
}
