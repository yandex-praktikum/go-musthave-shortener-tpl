package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spinel/go-musthave-shortener-tpl/internal/app/helper"
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/model"
	"github.com/spinel/go-musthave-shortener-tpl/internal/app/repository"
)

const Host = "http://localhost:8080"

func CreateShortenerHandler(repo *repository.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "wrong body", http.StatusBadRequest)
			return
		}
		url := string(body)
		if url == "" {
			http.Error(w, "no body", http.StatusBadRequest)
			return
		}
		shortener := &model.Shortener{
			URL: url,
		}
		var code helper.GeneratedString
		for {
			code, err = helper.NewGeneratedString()
			if err != nil {
				http.Error(w, "save shortener error", http.StatusInternalServerError)
				return
			}
			if !repo.Shortener.IncludesCode(string(code)) {
				break
			}
		}
		codeString := string(code)
		err = repo.Shortener.SaveShortener(codeString, shortener)
		if err != nil {
			http.Error(w, "shortener exists", http.StatusInternalServerError)
			return
		}
		result := fmt.Sprintf("%s/%s", Host, code)
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(result))
	}
}

func GetShortenerHandler(repo *repository.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		id := url[1:]
		if id == "" {
			http.Error(w, "no id", http.StatusBadRequest)
			return
		}
		shortener, err := repo.Shortener.GetShortenerBy(id)
		if err != nil {
			http.Error(w, "get shortener error", http.StatusInternalServerError)
			return
		}
		if shortener == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, shortener.URL, http.StatusTemporaryRedirect)
	}
}
