package main

import (
	"context"
	"log"
	"net/http"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/handler"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"

	_ "github.com/lib/pq"
)

func main() {

	config := configs.NewConfig()
	ctx := context.Background()
	db, err := repository.NewDBClient(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	repository.Migration(config)

	r := repository.NewStorage(db)
	s := service.NewService(r, config)
	h := handler.NewHandler(s)

	//start server
	httpServer := &http.Server{Addr: config.ServerAdress, Handler: h.InitRoutes()}
	httpServer.ListenAndServe()
}
