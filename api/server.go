package api

import (
	"net/http"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
)

type URLShortenerServer struct {
	http.Server
}

// New makes an instance of HTTP server that is ready to run
func New(
	shortenerSrv shortener.URLService,
	authSrv auth.IDService,
	addr string, baseURL url.URL,
) *URLShortenerServer {
	server := URLShortenerServer{
		Server: http.Server{
			Addr:    addr,
			Handler: handler.New(shortenerSrv, authSrv, baseURL),
		},
	}

	return &server
}
