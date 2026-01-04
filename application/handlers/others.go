package handlers

import (
	"net/http"
	"os"

	"github.com/rnium/rhttp/internal/headers"
	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/response"
)

func HealthCheck(r *request.Request) *response.Response {
	p := []byte("<h1>Everything working fine</h1>")
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return response.NewResponse(response.StatusOK, p, headers)
}

func Ping(r *request.Request) *response.Response {
	p := []byte("Pong")
	return response.NewResponse(response.StatusOK, p, nil)
}

func HttpBinStream(r *request.Request) *response.Response {
	n, _ := r.Param("n")
	res, err := http.Get("https://httpbin.org/stream/" + n)
	if err != nil {
		panic(err)
	}
	return response.NewChunkedResponse(response.StatusOK, res.Body, nil)
}

func Index(r *request.Request) *response.Response {
	f, err := os.Open("./templates/index.html")
	if err != nil {
		panic(err)
	}
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return response.NewChunkedResponse(response.StatusOK, f, headers)
}
