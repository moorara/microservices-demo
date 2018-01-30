package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "http://service/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, r)
	res := w.Result()

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestGetNotFoundHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "http://service/404", nil)
	w := httptest.NewRecorder()

	logger := log.NewNopLogger()
	notFoundHandler := GetNotFoundHandler(logger)
	notFoundHandler(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusNotFound)
	assert.Contains(t, string(body), `{"message":"Not found"}`)
}
