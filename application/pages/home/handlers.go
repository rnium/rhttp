package home

import (
	"os"

	"github.com/rnium/rhttp/internal/http/headers"
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
)

func healthCheck(r *request.Request) *response.Response {
	p := []byte("<h1>Everything working fine</h1>")
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return response.NewResponse(response.StatusOK, p, headers)
}

func index(r *request.Request) *response.Response {
	f, err := os.Open("./templates/index.html")
	if err != nil {
		panic(err)
	}
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return response.NewChunkedResponse(response.StatusOK, f, headers)
}
