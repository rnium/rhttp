package apis

import (
	"github.com/rnium/rhttp/internal/application/apis/cookies"
	"github.com/rnium/rhttp/internal/application/apis/dynamicdata"
	"github.com/rnium/rhttp/internal/application/apis/images"
	"github.com/rnium/rhttp/internal/application/apis/methods"
	"github.com/rnium/rhttp/internal/application/apis/request"
	"github.com/rnium/rhttp/internal/application/apis/response"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func RegisterRoutes(r *rhttp.Router) {
	methods.Register(r)
	request.Register(r)
	response.Register(r)
	dynamicdata.Register(r)
	images.Register(r)
	cookies.Register(r)
}
