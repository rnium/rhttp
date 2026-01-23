package home

import (
	"github.com/rnium/rhttp/internal/respond"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func healthCheck(r *rhttp.Request) *rhttp.Response {
	p := []byte("<h1>Everything working fine</h1>")
	res := rhttp.NewResponse(rhttp.StatusOK, p)
	_ = res.SetHeader("content-type", "text/html")
	return res
}

func index(r *rhttp.Request) *rhttp.Response {
	return respond.FileResponse("./web/templates/index.html")
}
