package apis

import (
	"github.com/rnium/rhttp/application/apis/methods"
	"github.com/rnium/rhttp/internal/router"
)


func RegisterRoutes(r *router.Router) {
	methods.Register(r)
}