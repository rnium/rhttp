package apis

import (
	"github.com/rnium/rhttp/application/apis/methods"
	requestinspection "github.com/rnium/rhttp/application/apis/request_inspection"
	"github.com/rnium/rhttp/internal/router"
)


func RegisterRoutes(r *router.Router) {
	methods.Register(r)
	requestinspection.Register(r)
}