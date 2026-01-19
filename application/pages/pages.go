package pages

import (
	"github.com/rnium/rhttp/application/pages/home"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func RegisterRoutes(r *rhttp.Router) {
	home.Register(r)
}
