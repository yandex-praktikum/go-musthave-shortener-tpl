package main

import (
	"net/http"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/handler"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/models/file"
)

func main() {
	var model file.Model

	config := configs.NewConfig()
	r := repository.NewStorage()
	s := service.NewService(r, &model, config)
	h := handler.NewHandler(s)

	//load data from file
	s.LoadFromFile()

	//start server
	httpServer := &http.Server{Addr: config.ServerAdress, Handler: h.InitRoutes()}
	httpServer.ListenAndServe()

}
