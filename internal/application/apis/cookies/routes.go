package cookies

import "github.com/rnium/rhttp/pkg/rhttp"

func Register(r *rhttp.Router) {
	r.Get("/cookies", cookiesHandler)
	r.Get("/cookies/set", setCookieHandler)
}