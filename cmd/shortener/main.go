package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/internal/app"
)

func main() {
	conf := app.LoadConfig()

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
