package application

import (
	"github.com/rnium/rhttp/internal/application/apis"
	"github.com/rnium/rhttp/internal/application/pages"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func getRoutes() *rhttp.Router {
	router := rhttp.NewRouter()
	pages.RegisterRoutes(router)
	apis.RegisterRoutes(router)
	return router
}
