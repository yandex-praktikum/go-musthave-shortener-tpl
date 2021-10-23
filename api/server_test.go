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
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/mocks"
	"github.com/stretchr/testify/require"
)

var baseURL = newURL("http://localhost:8080")
var longURL = newURL("http://test.com")

func TestHandlePostLongURL(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL.String()))
	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")

	urlService := new(urlmocks.URLServiceMock)
	urlService.On("ShortenURL", urlToShorten).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)

	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()

	h := handler.New(urlService, idService, pinger, baseURL)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, "http://localhost:8080/123", string(body))
}

func TestHandlePostLongURLConflict(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL.String()))
	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")

	urlService := new(urlmocks.URLServiceMock)
	urlService.On("ShortenURL", urlToShorten).Return(nil, model.ErrDuplicateURL)
	urlService.On("LookupURL", longURL).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)

	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()

	h := handler.New(urlService, idService, pinger, baseURL)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusConflict, res.StatusCode, "status code")
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
	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()
	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")
	h := handler.New(urlService, idService, pinger, baseURL)
	urlService.On("ShortenURL", urlToShorten).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "Content-Type")
	require.JSONEq(t, `{"result": "http://localhost:8080/123"}`, string(body))
}

func TestHandlePostApiShortenConflict(t *testing.T) {
	rw := httptest.NewRecorder()
	testLongURLJson := fmt.Sprintf(`{"url": "%s"}`, &longURL)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(testLongURLJson),
	)

	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")

	urlService := new(urlmocks.URLServiceMock)
	urlService.On("ShortenURL", urlToShorten).Return(nil, model.ErrDuplicateURL)
	urlService.On("LookupURL", longURL).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)

	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()

	h := handler.New(urlService, idService, pinger, baseURL)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusConflict, res.StatusCode, "status code")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "Content-Type")
	require.JSONEq(t, `{"result": "http://localhost:8080/123"}`, string(body))
}

func TestHandlePostShortenBatch(t *testing.T) {
	rw := httptest.NewRecorder()
	testLongURLBatchJSON := fmt.Sprintf(`[
			{
				"correlation_id": "abc",
				"original_url": "%s"
			}
		]`, &longURL)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		bytes.NewBufferString(testLongURLBatchJSON),
	)
	urlService := new(urlmocks.URLServiceMock)
	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()
	urlToShorten := model.NewURLToShorten(0, longURL)
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	absoluteShortURL, _ := url.Parse("http://localhost:8080/123")
	h := handler.New(urlService, idService, pinger, baseURL)
	urlService.On("ShortenURL", urlToShorten).Return(&shortenedURL, nil)
	urlService.On("AbsoluteURL", shortenedURL).Return(absoluteShortURL, nil)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "Content-Type")
	require.JSONEq(t, `[
			{
					"correlation_id": "abc",
					"short_url": "http://localhost:8080/123"
			}
		]`, string(body))
}

func TestHandleGetShortUrl(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	urlService := new(urlmocks.URLServiceMock)
	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()
	shortenedURL := model.NewShortenedURL(0, 123, longURL)
	h := handler.New(urlService, idService, pinger, baseURL)
	urlService.On("GetURLByID", 123).Return(&shortenedURL, nil)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "status code")
	require.Equal(t, longURL.String(), res.Header.Get("Location"), "location")
}

func TestHandleGetShortUrlNotFound(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	urlService := new(urlmocks.URLServiceMock)
	idService := authmocks.NewIDServiceStub()
	pinger := mocks.NewPingerStub()
	h := handler.New(urlService, idService, pinger, baseURL)
	urlService.On("GetURLByID", 123).Return(nil, model.ErrURLNotFound)

	h.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusNotFound, res.StatusCode, "status code")
}

func newURL(urlStr string) url.URL {
	u, errParse := url.Parse(urlStr)
	if errParse != nil {
		log.Fatalf("Cannot parse url [%s]: %s", urlStr, errParse.Error())
	}
	return *u
}
