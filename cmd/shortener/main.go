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

	config := configs.NewConfig()

	db, err := repository.NewDBClient(context.TODO(), config)
	if err != nil {
		log.Fatal(err)
	}

	repository.Migration(config)

	r := repository.NewStorage(db)
	s := service.NewService(r, config)
	h := handler.NewHandler(s)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//router.Use(handler.AuthMiddleware(h))
	router.GET("/:id", h.HandlerURLRelocation)
	router.GET("user/urls")
	router.GET("api/shorten/batch")
	router.GET("/ping", h.HandlerPingDB)
	router.POST("/", h.HandlerPost)
	router.POST("/api/shorten", h.HandlerPost)
	router.NoRoute(func(c *gin.Context) { c.String(http.StatusBadRequest, "Not allowed requset") })

	//start server
	router.Run(config.ServerAdress)

}
