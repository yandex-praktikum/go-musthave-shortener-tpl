package handler

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/middleware"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type URLShortenerHandler struct {
	*chi.Mux
	BaseURL url.URL
	Service shortener.URLService
	Pinger  storage.Pinger
}

func New(shortenerService shortener.URLService, idService auth.IDService, pinger storage.Pinger, baseURL url.URL) *URLShortenerHandler {
	h := &URLShortenerHandler{
		Mux:     chi.NewMux(),
		BaseURL: baseURL,
		Service: shortenerService,
		Pinger:  pinger,
	}

	h.Group(func(r chi.Router) {
		r.Use(middleware.Authenticator(idService))
		r.Use(middleware.GzipDecompressor)
		r.Use(middleware.GzipCompressor)
		r.Post("/", h.handlePostLongURL)
		r.Post("/api/shorten", h.handlePostAPIShorten)
		r.Post("/api/shorten/batch", h.handlePostShortenBatch)
		r.Get("/{id}", h.handleGetShortURL)
		r.Get("/user/urls", h.handleGetUserURLs)
	})

	h.Group(func(r chi.Router) {
		r.Get("/ping", h.handleGetPing)
	})

	return h
}

func userID(r *http.Request) int64 {
	return r.Context().Value(middleware.AuthContextKeyType{}).(int64)
}
