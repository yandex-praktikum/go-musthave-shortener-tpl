package handler

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}
type Request struct {
	LongURL string `json:"url"  binding:"required"`
	body    []byte
}
type Result struct {
	Result string `json:"result"`
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

//=================================================================
func (h *Handler) HandlerGet(c *gin.Context) {
	id := c.Param("id")
	longURL, err := h.service.GetURL(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", longURL)
}

//==================================================================
func (h *Handler) HandlerPostText(c *gin.Context) {
	var request Request
	//if request body is compressed
	if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			log.Fatal("error decoding response", err)
		}
		defer reader.Close()
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("error decoding response", err)
		}
		request.body = body

		//if request body is uncompressed
	} else {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "Not allowed request")
			return
		}
		request.body = body
	}

	if len(request.body) == 0 {
		c.String(http.StatusBadRequest, "Not allowed request")
		return
	}
	//Ganerate short URL and save to storage
	id, err := h.service.SaveURL(string(request.body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	}
	//if the client supports compression
	if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		gz.Write([]byte(fmt.Sprint(h.service.Config.BaseURL, "/", id)))
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Content-Type", "application/x-gzip")
		c.Status(http.StatusCreated)
		//if the client doesn't support compression
	} else {
		c.String(http.StatusCreated, fmt.Sprint(h.service.Config.BaseURL, "/", id))
	}
}

//===================================================================
func (h *Handler) HandlerPostJSON(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	if request.LongURL == "" || c.GetHeader("content-type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}

	id, err := h.service.SaveURL(request.LongURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	}
	shortURL := fmt.Sprint(h.service.Config.BaseURL, "/", id)
	var result Result
	result.Result = shortURL
	c.JSON(http.StatusCreated, result)
}
