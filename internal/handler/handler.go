package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
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
	longURL := h.service.GetURL(id)
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
	id := h.service.SaveURL(string(body))
	c.String(http.StatusCreated, "http://localhost:8080/"+id)

}
