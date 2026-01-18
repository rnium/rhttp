package home

import "github.com/rnium/rhttp/internal/router"


func Register(r *router.Router) {
	r.Get("/", index)
	r.Get("/health", healthCheck)
}