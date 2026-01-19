package methods

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/get", handleMethod)
	r.Post("/post", handleMethod)
	r.Put("/put", handleMethod)
	r.Patch("/patch", handleMethod)
	r.Delete("/delete", handleMethod)
}
