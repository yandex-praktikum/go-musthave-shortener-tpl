package main

import (
	"net/http"

	"github.com/EMus88/GO-Yandex-Study/internal/app/service"
	"github.com/EMus88/GO-Yandex-Study/internal/handler"
	"github.com/EMus88/GO-Yandex-Study/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	r := repository.NewStorage()
	s := service.NewService(r)
	h := handler.NewHandler(s)

	router := gin.Default()
	router.GET("/:id", h.HandlerGet)
	router.POST("/", h.HandlerPost)
	router.NoRoute(func(c *gin.Context) { c.String(http.StatusBadRequest, "Not allowed requset") })

	router.Run("localhost:8080")

}
