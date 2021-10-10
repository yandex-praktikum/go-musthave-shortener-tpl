package api

import (
	"net/http"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/inmem"
)

type URLShortenerServer struct {
	http.Server
}

// New makes an instance of HTTP server that is ready to run
func New(
	shortenerSrv shortener.URLShortener,
	authSrv auth.IDService,
	addr string, baseURL url.URL,
) *URLShortenerServer {
	storage := inmem.New()
	server := URLShortenerServer{
		Server: http.Server{
			Addr:    addr,
			Handler: handler.New(storage, authSrv, baseURL),
		},
	}

	return &server
}
