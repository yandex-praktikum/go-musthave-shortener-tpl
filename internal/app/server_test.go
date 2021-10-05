package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/mocks"
	"github.com/stretchr/testify/require"
)

var BaseURL url.URL

func TestMain(m *testing.M) {
	url, errUrl := url.Parse("http://localhost:8080")
	if errUrl != nil {
		log.Fatalf("Cannot parse base URL: %s", errUrl.Error())
	}

	BaseURL = *url

	code := m.Run()
	os.Exit(code)
}

func TestHandlePostLongURL(t *testing.T) {
	rw := httptest.NewRecorder()
	testRawURL := "http://test.com"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testRawURL))
	repo := new(mocks.StorageMock)
	url, _ := url.Parse(testRawURL)
	storableURL := model.NewStorableURL(url)
	storeURL := model.NewStoreURL(123, url)
	sh := NewURLShortener(repo, BaseURL)
	repo.On("Save", storableURL).Return(storeURL)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, "http://localhost:8080/123", string(body))
}

func TestHandlePostApiShorten(t *testing.T) {
	rw := httptest.NewRecorder()
	testRawLongURL := "http://test.com"
	testLongURLJson := fmt.Sprintf(`{"url": "%s"}`, testRawLongURL)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(testLongURLJson),
	)
	repo := new(mocks.StorageMock)
	url, _ := url.Parse(testRawLongURL)
	storableURL := model.NewStorableURL(url)
	storeURL := model.NewStoreURL(123, url)
	sh := NewURLShortener(repo, BaseURL)
	repo.On("Save", storableURL).Return(storeURL)

	sh.ServeHTTP(rw, req)

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
	repo := new(mocks.StorageMock)
	testRawURL := "http://test.com"
	url, _ := url.Parse(testRawURL)
	storeURL := model.NewStoreURL(123, url)
	sh := NewURLShortener(repo, BaseURL)
	repo.On("GetByID", 123).Return(&storeURL)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "status code")
	require.Equal(t, testRawURL, res.Header.Get("Location"), "location")
}
