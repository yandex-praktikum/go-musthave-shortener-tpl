package handler

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/middleware"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
)

type URLShortenerHandler struct {
	*chi.Mux
	BaseURL url.URL
	Service shortener.URLService
}

func New(shortenerService shortener.URLService, idService auth.IDService, baseURL url.URL) *URLShortenerHandler {
	h := &URLShortenerHandler{
		Mux:     chi.NewMux(),
		BaseURL: baseURL,
		Service: shortenerService,
	}

	auth := middleware.Authenticator{IDService: idService}
	h.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)
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

func userID(r *http.Request) int {
	return r.Context().Value(middleware.AuthContextKeyType{}).(int)
}
