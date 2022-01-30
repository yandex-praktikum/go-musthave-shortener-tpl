package handler

import (
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/EMus88/go-musthave-shortener-tpl/internal/repository/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *service.Service
	publicKey string
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
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Content-Type", "application/x-gzip")

	} else {
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
		sessionID, cookieErr := c.Cookie("session")
		//if cookie is empty
		if cookieErr != nil {
			if !errors.Is(cookieErr, http.ErrNoCookie) {
				c.String(http.StatusInternalServerError, "Auth error")
				log.Fatal(cookieErr)
			}
			//create new cookie
			encID, err := h.service.CreateNewSession()

			if err != nil {
				c.String(http.StatusInternalServerError, "Auth error")
				log.Fatal(err)
			}
			sessionID = encID
			//if cookie is not empty
		} else {
			//validation seesion value
			isValid := true
			byteID, DecErr := hex.DecodeString(sessionID)
			if DecErr != nil || len(byteID) != 16 {
				isValid = false
			}
			//if cookie value is valid
			if isValid {
				//read non-public key
				key, err := h.service.Auth.ReadSessionID(sessionID)
				if err != nil {
					//create new cookie
					encID, err := h.service.CreateNewSession()
					if err != nil {
						c.String(http.StatusInternalServerError, "Auth error")
						log.Fatal(err)
					}
					sessionID = encID
				} else {
					//try find non-public key in data base
					isFound := h.service.Repository.GetCookie(key)
					//if non-public key was not found
					if !isFound {
						//create new cookie
						encID, err := h.service.CreateNewSession()
						if err != nil {
							c.String(http.StatusInternalServerError, "Auth error")
							log.Fatal(err)
						}
						sessionID = encID

					}
				}
				//if cookie value is not valid
			} else {
				//create new cookie
				encID, err := h.service.CreateNewSession()
				if err != nil {
					c.String(http.StatusInternalServerError, "Auth error")
					log.Fatal(err)
				}
				sessionID = encID
			}

		}
		h.publicKey = sessionID
		c.SetCookie("session", sessionID, 3600, "", "localhost", false, true)

		c.Next()
	}
}

//=================================================================
func (h *Handler) HandlerURLRelocation(c *gin.Context) {
	id := c.Param("id")
	longURL, err := h.service.GetURL(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, longURL)
}

//=================================================================
func (h *Handler) HandlerGetList(c *gin.Context) {

	key, _ := h.service.Auth.ReadSessionID(h.publicKey)
	list, err := h.service.Repository.GetList(key)
	if err != nil {
		c.String(http.StatusNoContent, "Not found data")
		return
	}
	data, err := json.Marshal(list)
	if err != nil {
		c.String(http.StatusNoContent, "Not found data")
	}
	c.Data(http.StatusOK, gin.MIMEJSON, data)

}

//=================================================================
func (h *Handler) HandlerSaveBatch(c *gin.Context) {
	var model []model.BatchRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.String(http.StatusBadRequest, "Not allowed request")
		return
	}
	list, err := h.service.SaveBatch(model, h.publicKey)

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal error")
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		c.String(http.StatusNoContent, "Not found data")
		return
	}
	c.Data(http.StatusCreated, gin.MIMEJSON, data)
}

//=================================================================
func (h *Handler) HandlerPingDB(c *gin.Context) {
	if err := h.service.Repository.PingDB(); err != nil {
		c.String(http.StatusInternalServerError, "DB connection is not available")
	}
	c.String(http.StatusOK, "DB connection succes")
}

//==================================================================
func (h *Handler) HandlerPostURL(c *gin.Context) {
	if c.Request.URL.Path == "/api/shorten" {
		c.Header("content-type", "application/json")
	}
	request, err := parseRequest(c)
	if err != nil {
		c.String(http.StatusBadRequest, "error: Not Allowd request")
		return
	}
	shortURL, err := h.service.SaveURL(string(request.body), h.publicKey)
	if err != nil {
		if errors.Is(err, errors.Unwrap(err)) {
			c.String(http.StatusConflict, "")
			var response Response
			response.ShortURL = shortURL
			renderResponse(c, &response)
		} else {
			c.String(http.StatusInternalServerError, "error: Internal error")
		}
	} else {
		var response Response
		response.ShortURL = shortURL
		renderResponse(c, &response)
	}
}
