package images

import (
	"github.com/rnium/rhttp/internal/respond"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func imageHandler(r *rhttp.Request) *rhttp.Response {
	accept, _ := r.Headers.Get("accept")
	path := getImagePath(accept)
	return respond.FileResponse(path)
}

func handlerFactory(accept string) rhttp.Handler {
	return func(r *rhttp.Request) *rhttp.Response {
		_, _ = r.Headers.Replace("accept", accept)
		return imageHandler(r)
	}
}
