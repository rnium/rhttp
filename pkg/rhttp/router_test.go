package rhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func demoHandler(request *Request) *Response {
	p := []byte("hello world")
	return NewResponse(StatusOK, p, nil)
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
