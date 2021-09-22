package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var BaseUrl url.URL

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetURLBy(id int) *url.URL {
	args := r.Called(id)
	return args.Get(0).(*url.URL)
}

func (r *RepositoryMock) SaveURL(u url.URL) int {
	args := r.Called(u)
	return args.Int(0)
}

func TestMain(m *testing.M) {
	url, errConf := url.Parse("http://localhost:8080")
	if errConf != nil {
		panic(errConf)
	}

	BaseUrl = *url

	code := m.Run()
	os.Exit(code)
}

func TestHandlePostLongURL(t *testing.T) {
	rw := httptest.NewRecorder()
	testRawURL := "http://test.com"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testRawURL))
	repo := new(RepositoryMock)
	url, _ := url.Parse(testRawURL)
	sh := NewURLShortener(repo, BaseUrl)
	repo.On("SaveURL", *url).Return(123)

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
	repo := new(RepositoryMock)
	url, _ := url.Parse(testRawLongURL)
	sh := NewURLShortener(repo, BaseUrl)
	repo.On("SaveURL", *url).Return(123)

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
	repo := new(RepositoryMock)
	testRawURL := "http://test.com"
	url, _ := url.Parse(testRawURL)
	sh := NewURLShortener(repo, BaseUrl)
	repo.On("GetURLBy", 123).Return(url)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "status code")
	require.Equal(t, testRawURL, res.Header.Get("Location"), "location")
}
