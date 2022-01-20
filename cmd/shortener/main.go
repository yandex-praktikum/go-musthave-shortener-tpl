package main

import (
	"net/http"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/handler"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/models/file"
	"github.com/gin-gonic/gin"
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/:id", h.HandlerGet)
	router.POST("/", h.HandlerPostText)
	router.POST("/api/shorten", h.HandlerPostJSON)
	router.NoRoute(func(c *gin.Context) { c.String(http.StatusBadRequest, "Not allowed requset") })

	router.Run(config.ServerAdress)
}
