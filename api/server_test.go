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
	storagemocks "github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/mocks"
	"github.com/stretchr/testify/require"
)

var baseURL = newURL("http://localhost:8080")
var longURL = newURL("http://test.com")

func TestHandlePostLongURL(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(longURL.String()))
	storage := new(storagemocks.StorageMock)
	idService := new(authmocks.IDServiceMock)
	storableURL := model.NewURLToShorten(0, longURL)
	storeURL := model.NewShortenedURL(0, 123, longURL)
	h := handler.New(storage, idService, baseURL)
	storage.On("Save", storableURL).Return(storeURL)
	idService.On("SignUp").Return(&model.User{ID: 0, Key: []byte("")}, nil)

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
	storage := new(storagemocks.StorageMock)
	idService := new(authmocks.IDServiceMock)
	storableURL := model.NewURLToShorten(0, longURL)
	storeURL := model.NewShortenedURL(0, 123, longURL)
	h := handler.New(storage, idService, baseURL)
	storage.On("Save", storableURL).Return(storeURL)
	idService.On("SignUp").Return(&model.User{ID: 0, Key: []byte("")}, nil)

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
	storage := new(storagemocks.StorageMock)
	idService := new(authmocks.IDServiceMock)
	storeURL := model.NewShortenedURL(0, 123, longURL)
	h := handler.New(storage, idService, baseURL)
	storage.On("GetByID", 123).Return(&storeURL)
	idService.On("SignUp").Return(&model.User{ID: 0, Key: []byte("")}, nil)

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
