package main

import (
	"net/http"

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

func MyBad(r *request.Request) *response.Response {
	panic("Oops.. My Bad")
}

func HttpBinStream(r *request.Request) *response.Response {
	n, _ := r.Param("n")
	res, err := http.Get("https://httpbin.org/stream/" + n)
	if err != nil {
		panic(err)
	}
	return response.NewChunkedResponse(response.StatusOK, res.Body, nil)
}
