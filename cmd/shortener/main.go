package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/config"
)

func main() {
	conf, errConf := config.Load()
	if errConf != nil {
		log.Fatalf("Cannot load config: %s", errConf.Error())
	}

	server := api.New(*conf)
	log.Println("Starting server...")

	go start(server)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	if errShutdown := server.Shutdown(context.Background()); errShutdown != nil {
		log.Fatalf("Could not gracefully stop the server: %s", errShutdown.Error())
	}
	log.Println("Server stopped.")
}

func start(s *api.URLShortenerServer) {
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Cannot start the server: %v", err.Error())
	}
}
