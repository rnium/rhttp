package home

import (
	"os"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func healthCheck(r *rhttp.Request) *rhttp.Response {
	p := []byte("<h1>Everything working fine</h1>")
	headers := rhttp.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return rhttp.NewResponse(rhttp.StatusOK, p, headers)
}

func index(r *rhttp.Request) *rhttp.Response {
	f, err := os.Open("./templates/index.html")
	if err != nil {
		panic(err)
	}
	headers := rhttp.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return rhttp.NewChunkedResponse(rhttp.StatusOK, f, headers)
}
