package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/internal/app"
)

func main() {
	conf, errConf := app.LoadConfig()
	if errConf != nil {
		log.Fatalf("Cannot load config: %s", errConf.Error())
	}

	server := app.NewServer(conf)
	log.Println("Starting server...")

	go start(server)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	server.Shutdown(context.Background())
	log.Println("Server stopped.")
}

func start(s *app.URLShortenerServer) {
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Cannot start the server: %v", err.Error())
	}
}
