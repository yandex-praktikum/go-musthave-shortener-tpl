package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/inmem"
)

type URLShortenerServer struct {
	http.Server
	Storage storage.BulkStorage
}

// New makes an instance of HTTP server that is ready to run
func New(shortenerSrv shortener.URLShortener, authSrv auth.IDService, addr string, baseURL url.URL) *URLShortenerServer {
	storage := inmem.New()
	server := URLShortenerServer{
		Server: http.Server{
			Addr:    addr,
			Handler: handler.New(storage, authSrv, baseURL),
		},
		Storage: storage,
	}

	return &server
}

// ListenAndServe restores the server state from the backup file
// and starts the HTTP server
func (s *URLShortenerServer) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

// Shutdown backs-up the server state into the backup file
// and stops the HTTP server gracefully
func (s *URLShortenerServer) Shutdown(ctx context.Context) error {
	return nil
}
