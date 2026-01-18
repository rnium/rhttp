package application

import (
	"github.com/rnium/rhttp/application/apis"
	"github.com/rnium/rhttp/application/pages"
	"github.com/rnium/rhttp/internal/router"
)

func getRoutes() *router.Router {
	router := router.NewRouter()
	pages.RegisterRoutes(router)
	apis.RegisterRoutes(router)
	return router
}
