package main

import (
	"context"
	"log"
	"net/http"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/handler"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	//set new configuraion
	config := configs.NewConfig()

	//init data base connection
	db, err := repository.NewDBClient(context.TODO(), config)
	if err != nil {
		log.Fatal(err)
	}
	//creat tables in data bese by migration models
	repository.Migration(config)

	//init main modelus
	r := repository.NewStorage(db)
	s := service.NewService(r, config)
	h := handler.NewHandler(s)

	//set server settings
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(handler.AuthMiddleware(h))

	router.GET("/:id", h.HandlerURLRelocation)
	router.GET("user/urls", h.HandlerGetList)
	router.GET("/ping", h.HandlerPingDB)

	router.POST("/", h.HandlerPostURL)
	router.POST("/api/shorten", h.HandlerPostURL)
	router.POST("/api/shorten/batch", h.HandlerSaveBatch)

	router.NoRoute(func(c *gin.Context) { c.String(http.StatusBadRequest, "Not allowed requset") })

	//start server
	router.Run(config.ServerAdress)

}
