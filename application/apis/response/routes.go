package response

import "github.com/rnium/rhttp/internal/router"

func Register(r *router.Router) {
	r.Get("/cache", cache)
}