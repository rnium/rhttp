package application

import (
	"github.com/rnium/rhttp/application/handlers"
	"github.com/rnium/rhttp/internal/router"
)

func getRoutes() *router.Router {
	router := router.NewRouter()
	router.Post("/post", handlers.HandlePost)
	router.Get("/", handlers.Index)
	return router
}
