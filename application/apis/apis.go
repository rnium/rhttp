package apis

import (
	"github.com/rnium/rhttp/application/apis/methods"
	"github.com/rnium/rhttp/application/apis/request"
	"github.com/rnium/rhttp/application/apis/response"
	"github.com/rnium/rhttp/internal/router"
)


func RegisterRoutes(r *router.Router) {
	methods.Register(r)
	request.Register(r)
	response.Register(r)
}