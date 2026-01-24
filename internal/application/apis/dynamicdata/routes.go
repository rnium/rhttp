package dynamicdata

import "github.com/rnium/rhttp/pkg/rhttp"


func Register(r *rhttp.Router) {
	r.Get("/base64/:value", decodeBase64Handler)
	r.Get("/uuid", uuidGenHandler)
}