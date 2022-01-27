package handler

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *service.Service
	sessionID string
}
type Request struct {
	LongURL string `json:"url"  binding:"required"`
	body    []byte
}
type Response struct {
	ShortURL string `json:"result"`
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

//================================================================
func isEncodingSupport(c *gin.Context) bool {
	//if the client supports compression
	if strings.Contains(c.GetHeader("Accept-Encoding"), "qweqwe") {
		return true
		//if the client doesn't support compression
	} else {
		return false
	}
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
		request.body = []byte(request.LongURL)

	case "text/plain":
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		if len(body) < 1 {
			return nil, errors.New("error")
		}
		if ok := govalidator.IsURL(string(body)); !ok {
			return nil, errors.New("error")
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
	}
	return &request, nil
}

//================================================================
func renderResponse(c *gin.Context, response *Response) {

	if isEncodingSupport(c) {
		c.Status(http.StatusCreated)
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		gz.Write([]byte(response.ShortURL))
		c.Writer.Header().Set("Content-Encoding", "sdfsdf")
		c.Writer.Header().Set("Content-Type", "application/x-gzip")

	} else {
		fmt.Println(c.Request.Header.Get("content-type"))
		switch c.Request.Header.Get("content-type") {
		case "application/json":
			c.JSON(http.StatusCreated, response)
		case "text/plain":
			c.String(http.StatusCreated, response.ShortURL)
		default:
			c.String(http.StatusCreated, response.ShortURL)
		}
	}
}

//==================================================

func AuthMiddleware(h *Handler) gin.HandlerFunc {

	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session")

		//validation session value
		//var id int = h.service.Repository.GetCookieID("sdfsdfsdf")
		//fmt.Println(id)

		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Fatal(err)
			}
			_, encID, err := h.service.Auth.CreateSissionID()
			if err != nil {
				log.Fatal(err)
			}
			sessionID = encID
		}

		c.SetCookie("session", sessionID, 3600, "", "localhost", false, true)

		c.Next()
	}
}

//=================================================================
func (h *Handler) HandlerURLRelocation(c *gin.Context) {
	log.Println("GET hanler")
	//id := c.Param("id")
	//longURL, err := h.service.GetURL(id)
	longURL := "https://github.com/EMus88/go-musthave-shortener-tpl"
	// if err != nil {
	// 	c.String(http.StatusBadRequest, err.Error())
	// 	return
	// }
	c.Redirect(http.StatusTemporaryRedirect, longURL)
}

//=================================================================
func (h *Handler) HandlerPingDB(c *gin.Context) {
	if err := h.service.Repository.PingDB(); err != nil {
		c.String(http.StatusInternalServerError, "DB connection is not available")
	}
	c.String(http.StatusOK, "DB connection succes")
}

//==================================================================
func (h *Handler) HandlerPost(c *gin.Context) {
	request, err := parseRequest(c)
	if err != nil {
		c.String(http.StatusBadRequest, "error: Not Allowd request")
		return
	}

	shortURL, err := h.service.SaveURL(string(request.body), h.sessionID)
	if err != nil {
		c.String(http.StatusInternalServerError, "error: Internal error")
	}
	//write result
	var response Response
	response.ShortURL = shortURL
	renderResponse(c, &response)
}
