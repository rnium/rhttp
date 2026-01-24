package response

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/cache", cache)
	r.Get("/cache/:value", setCacheCtrl)
	r.Get("/etag/:etag", etagHandler)
	r.Get("/response-headers", responseHeaders)
	r.Post("/response-headers", responseHeaders)
}
