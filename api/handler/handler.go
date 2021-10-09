package handler

import (
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/middleware"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
	shortenerV1 "github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener/v1"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type URLShortenerHandler struct {
	*chi.Mux
	Storage storage.Storage
	BaseURL url.URL
	Service shortener.URLShortener
}

func New(s storage.Storage, baseURL url.URL) *URLShortenerHandler {
	h := &URLShortenerHandler{
		Mux:     chi.NewMux(),
		Storage: s,
		BaseURL: baseURL,
		Service: shortenerV1.New(s, baseURL),
	}
	h.Use(middleware.GzipDecompressor)
	h.Use(middleware.GzipCompressor)
	h.Post("/", h.handlePostLongURL)
	h.Post("/api/shorten", h.handlePostAPIShorten)
	h.Get("/{id}", h.handleGetShortURL)

	return h
}
