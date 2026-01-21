package home

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/", index)
	r.Get("/health", healthCheck)
}
