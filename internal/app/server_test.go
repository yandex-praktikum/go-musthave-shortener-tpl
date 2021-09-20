package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetURLBy(id int) *url.URL {
	args := r.Called(id)
	return args.Get(0).(*url.URL)
}

func (r *RepositoryMock) SaveURL(u *url.URL) int {
	args := r.Called(u)
	return args.Int(0)
}

func TestHandlePostLongURL(t *testing.T) {
	rw := httptest.NewRecorder()
	testRawURL := "http://test.com"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testRawURL))
	repo := new(RepositoryMock)
	url, _ := url.Parse(testRawURL)
	sh := NewURLShortener(repo)
	repo.On("SaveURL", url).Return(123)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, string(body), "http://localhost:8080/123")
}

func TestHandleGetShortUrl(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	repo := new(RepositoryMock)
	testRawURL := "http://test.com"
	url, _ := url.Parse(testRawURL)
	sh := NewURLShortener(repo)
	repo.On("GetURLBy", 123).Return(url)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "status code")
	require.Equal(t, testRawURL, res.Header.Get("Location"), "location")
}
