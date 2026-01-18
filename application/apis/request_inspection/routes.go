package requestinspection

import "github.com/rnium/rhttp/internal/router"

func Register(r *router.Router) {
	r.Get("/headers", viewHeaders)
	r.Get("/ip", viewIp)
	r.Get("/user-agent", viewUserAgent)
}