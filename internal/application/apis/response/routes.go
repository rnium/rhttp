package response

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/cache", cache)
	r.Get("/cache/:value", setCacheCtrl)
}
