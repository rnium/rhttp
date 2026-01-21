package request

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/headers", viewHeaders)
	r.Get("/ip", viewIp)
	r.Get("/user-agent", viewUserAgent)
}
