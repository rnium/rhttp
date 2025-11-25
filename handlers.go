package main

import (
	"fmt"

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

func MethodChecker(r *request.Request) *response.Response {
	p := fmt.Appendf(nil, "Successfull handled %s request from methodchecker endpoint", r.RequestLine.Method)
	return response.NewResponse(response.StatusOK, p, nil)
}

func ParamChecker(r *request.Request) *response.Response {
	p_name := "pk"
	pk, _ := r.Param(p_name)
	p := fmt.Appendf(nil, "Parameter '%s' is: %s", p_name, pk)
	return response.NewResponse(response.StatusOK, p, nil)
}
