package pages

import (
	"github.com/rnium/rhttp/internal/application/pages/home"
	"github.com/rnium/rhttp/internal/application/pages/static"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func RegisterRoutes(r *rhttp.Router) {
	home.Register(r)
	static.Register(r)
}
