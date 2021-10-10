package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	authmocks "github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth/mocks"
	urlmocks "github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/shortener/mocks"
	"github.com/stretchr/testify/require"
)

var baseURL = newURL("http://localhost:8080")
var longURL = newURL("http://test.com")

func TestHandlePostLongURL(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL.String()))
	urlService := new(urlmocks.URLServiceMock)
	idService := new(authmocks.IDServiceMock)
	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")
	h := handler.New(urlService, idService, baseURL)
	urlService.On("ShortenURL", urlToShorten).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)
	idService.On("SignUp").Return(&model.User{ID: 0, Key: ""}, nil)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, "http://localhost:8080/123", string(body))
}

func TestHandlePostApiShorten(t *testing.T) {
	rw := httptest.NewRecorder()
	testLongURLJson := fmt.Sprintf(`{"url": "%s"}`, &longURL)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(testLongURLJson),
	)
	urlService := new(urlmocks.URLServiceMock)
	idService := new(authmocks.IDServiceMock)
	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")
	h := handler.New(urlService, idService, baseURL)
	urlService.On("ShortenURL", urlToShorten).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)
	idService.On("SignUp").Return(&model.User{ID: 0, Key: ""}, nil)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "Content-Type")
	require.JSONEq(t, `{"result": "http://localhost:8080/123"}`, string(body))
}

func TestHandleGetShortUrl(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	urlService := new(urlmocks.URLServiceMock)
	idService := new(authmocks.IDServiceMock)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	h := handler.New(urlService, idService, baseURL)
	urlService.On("GetByID", 123).Return(&shortenedURL, nil)
	idService.On("SignUp").Return(&model.User{ID: 0, Key: ""}, nil)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "status code")
	require.Equal(t, longURL.String(), res.Header.Get("Location"), "location")
}

func newURL(urlStr string) url.URL {
	u, errParse := url.Parse(urlStr)
	if errParse != nil {
		log.Fatalf("Cannot parse url [%s]: %s", urlStr, errParse.Error())
	}
	return *u
}
