package home

import (
	"os"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func healthCheck(r *rhttp.Request) *rhttp.Response {
	p := []byte("<h1>Everything working fine</h1>")
	res := rhttp.NewResponse(rhttp.StatusOK, p)
	_ = res.SetHeader("content-type", "text/html")
	return res
}

func index(r *rhttp.Request) *rhttp.Response {
	f, err := os.Open("./web/templates/index.html")
	if err != nil {
		panic(err)
	}

	res := rhttp.NewChunkedResponse(rhttp.StatusOK, f)
	_ = res.SetHeader("Content-Type", "text/html")
	return res
}
