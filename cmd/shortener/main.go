package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v6"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/internal/app"
)

func main() {
	var conf app.Config
	errConf := env.Parse(&conf)
	if errConf != nil {
		panic(errConf)
	}

	server := app.NewServer(conf)
	log.Println("Starting server...")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	server.Shutdown(context.Background())
	log.Println("Server stopped.")
}
