package static

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/static/*", webStatic)
}