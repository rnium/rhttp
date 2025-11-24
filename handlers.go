package main

import (
	"fmt"

	"github.com/rnium/rhttp/internal/headers"
	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/response"
)

func healthCheck(r *request.Request) *response.Response {
	p := []byte("<h1>Everything working fine</h1>")
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return response.NewResponse(response.StatusOK, p, headers)
}

func myBad(r *request.Request) *response.Response {
	panic("My Bad")
}

func methodChecker(r *request.Request) *response.Response {
	p := fmt.Appendf(nil, "Successfull handled %s request from methodchecker endpoint", r.RequestLine.Method)
	return response.NewResponse(response.StatusOK, p, nil)
}
