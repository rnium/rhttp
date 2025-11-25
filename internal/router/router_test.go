package router

import (
	"testing"

	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/response"
	"github.com/stretchr/testify/assert"
)

func demoHandler(request *request.Request) *response.Response {
	p := []byte("hello world")
	return response.NewResponse(response.StatusOK, p, nil)
}

func TestRouter(t *testing.T) {
	router := NewRouter()
	view, _ := router.getView("/hello")
	assert.Nil(t, view)
	router.Get("/hello", demoHandler)
	view, _ = router.getView("/hello")
	assert.NotNil(t, view)
	assert.Equal(t, []string{MethodGet}, view.methods)
	router.Post("/hello", demoHandler)
	assert.Equal(t, []string{MethodGet, MethodPost}, view.methods)
}
