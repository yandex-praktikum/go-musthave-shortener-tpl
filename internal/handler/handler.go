package handler

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/asaskevich/govalidator"
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
func parseRequest(c *gin.Context) (*Request, error) {
	var request Request
	switch c.Request.Header.Get("content-type") {

	case "application/json":
		if err := c.ShouldBindJSON(&request); err != nil {
			return nil, err
		}
		if ok := govalidator.IsURL(string(request.LongURL)); !ok {
			return nil, errors.New("error")
		}

	case "text/plain":
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		if len(body) < 1 {
			return nil, errors.New("error")
		}
		if ok := govalidator.IsURL(string(body)); !ok {
			return nil, err
		}
		request.body = body

	case "application/x-gzip":
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				return nil, err
			}
			defer reader.Close()
			body, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}
			request.body = body
		}
	default:
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		request.body = body
		fmt.Println("default")
	}
	fmt.Println(string(request.body))
	return &request, nil
}

//=================================================================
func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use()
	r.GET("/:id", h.HandlerGet)
	r.POST("/", h.HandlerPostText)
	r.POST("/api/shorten", h.HandlerPostJSON)
	r.NoRoute(func(c *gin.Context) { c.String(http.StatusBadRequest, "Not allowed requset") })
	return r
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
	request, err := parseRequest(c)
	if err != nil {
		c.String(http.StatusBadRequest, "error: Not Allowd request")
	}

	//Ganerate short URL and save to storage
	id, err := h.service.SaveURL(string(request.body))
	if err != nil {
		c.String(http.StatusInternalServerError, "error: Internal error")
	}

	//if the client supports compression
	if strings.Contains(c.GetHeader("Accept-Encoding"), "dfgdfg") {
		c.Status(http.StatusCreated)
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		gz.Write([]byte(fmt.Sprint(h.service.Config.BaseURL, "/", id)))
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Content-Type", "application/x-gzip")
		//if the client doesn't support compression
	} else {
		c.String(http.StatusCreated, fmt.Sprint(h.service.Config.BaseURL, "/", id))
	}
}

//===================================================================
func (h *Handler) HandlerPostJSON(c *gin.Context) {
	request, err := parseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not Allowd request"})
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
