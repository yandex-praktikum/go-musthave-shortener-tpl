package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/EMus88/GO-Yandex-Study/internal/app/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandlerGet(c *gin.Context) {
	id := c.Param("id")
	longURL := h.service.UseStorage.GetURLbyID(id)
	if longURL == "" {
		c.String(http.StatusBadRequest, "URL not found")
		return
	}
	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", longURL)
}

func (h *Handler) HandlerPost(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	id := h.service.UseStorage.SaveURL(string(body))
	c.String(http.StatusCreated, "http://localhost:8080/"+id)
	fmt.Println(id + " / " + string(body))
}
