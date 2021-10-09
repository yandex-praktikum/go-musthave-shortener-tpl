package handler

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/middleware"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
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

func New(s storage.Storage, idService auth.IDService, baseURL url.URL) *URLShortenerHandler {
	h := &URLShortenerHandler{
		Mux:     chi.NewMux(),
		Storage: s,
		BaseURL: baseURL,
		Service: shortenerV1.New(s, baseURL),
	}
	auth := middleware.Authenticator{IDService: idService}
	h.Use(auth.Authenticate)
	h.Use(middleware.GzipDecompressor)
	h.Use(middleware.GzipCompressor)
	h.Post("/", h.handlePostLongURL)
	h.Post("/api/shorten", h.handlePostAPIShorten)
	h.Get("/{id}", h.handleGetShortURL)
	h.Get("/user/urls", h.handleGetUserURLs)

	return h
}

func userID(r *http.Request) int {
	return r.Context().Value(middleware.AuthContextKeyType{}).(int)
}
