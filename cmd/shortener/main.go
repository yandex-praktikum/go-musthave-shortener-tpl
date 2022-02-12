package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/EMus88/go-musthave-shortener-tpl/configs"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/handler"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var done = make(chan os.Signal, 1)

func main() {

	//set new configuraion
	config := configs.NewConfig()

	//data base connection
	db, err := repository.NewDBClient(context.TODO(), config)
	if err != nil {
		log.Fatal(err)
	}
	//migration of database tables
	repository.Migration(config)
	if err != nil {
		log.Fatal(err)
	}

	//initialization of main modules
	r := repository.NewStorage(db)
	s := service.NewService(r, config)
	h := handler.NewHandler(s)

	go r.DeleteBufferRefreshing()
	//initializing router
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
	router.DELETE("/api/user/urls", h.HandlerDeleteURLs)

	router.NoRoute(func(c *gin.Context) { c.String(http.StatusBadRequest, "Not allowed requset") })

	log.Println("Routes inited")

	//start server
	go router.Run(config.ServerAdress)
	log.Printf("Server started: %s", config.ServerAdress)

	//shutdown
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	r.ShutdownChan <- struct{}{}
	<-r.DeleteCompleted
	log.Println("Server stopped")
}
