package api

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type URLShortenerServer struct {
	http.Server
}

// New makes an instance of HTTP server that is ready to run
func New(
	shortenerSrv shortener.URLService,
	authSrv auth.IDService,
	pinger storage.Pinger,
	addr string, baseURL url.URL,
) *URLShortenerServer {
	server := URLShortenerServer{
		Server: http.Server{
			Addr:    addr,
			Handler: handler.New(shortenerSrv, authSrv, pinger, baseURL),
		},
	}

	log.Println("Starting server...")

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Server failed: %s", err.Error())
		}
	}()

	return &server
}

func (s *URLShortenerServer) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped.")

	return nil
}
