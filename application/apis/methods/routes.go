package methods

import "github.com/rnium/rhttp/internal/router"

func Register(r *router.Router) {
	r.Get("/get", handleMethod)
	r.Post("/post", handleMethod)
	r.Put("/put", handleMethod)
	r.Patch("/patch", handleMethod)
	r.Delete("/delete", handleMethod)
}