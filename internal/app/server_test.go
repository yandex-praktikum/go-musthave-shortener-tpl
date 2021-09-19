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

func (r *RepositoryMock) GetUrlBy(id int) *url.URL {
	args := r.Called(id)
	return args.Get(0).(*url.URL)
}

func (r *RepositoryMock) SaveUrl(u *url.URL) int {
	args := r.Called(u)
	return args.Int(0)
}

func TestHandlePostLongUrl(t *testing.T) {
	rw := httptest.NewRecorder()
	testRawUrl := "http://test.com"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testRawUrl))
	repo := new(RepositoryMock)
	url, _ := url.Parse(testRawUrl)
	sh := NewUrlShortener(repo)
	repo.On("SaveUrl", url).Return(123)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	body, errBody := ioutil.ReadAll(res.Body)
	require.NoError(t, errBody)

	require.Equal(t, http.StatusCreated, res.StatusCode, "status code")
	require.Equal(t, string(body), "http://localhost:8080/123")
}

func TestHandleGetShortUrl(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	repo := new(RepositoryMock)
	testRawUrl := "http://test.com"
	url, _ := url.Parse(testRawUrl)
	sh := NewUrlShortener(repo)
	repo.On("GetUrlBy", 123).Return(url)

	sh.ServeHTTP(rw, req)

	res := rw.Result()
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "status code")
	require.Equal(t, testRawUrl, res.Header.Get("Location"), "location")
}
