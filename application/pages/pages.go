package pages

import (
	"github.com/rnium/rhttp/application/pages/home"
	"github.com/rnium/rhttp/internal/router"
)


func RegisterRoutes(r *router.Router) {
	home.Register(r)
}