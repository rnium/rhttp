package application

import (
	"github.com/rnium/rhttp/application/handlers"
	"github.com/rnium/rhttp/internal/router"
)

func getRoutes() *router.Router {
	router := router.NewRouter()
	router.Get("/", handlers.Index)
	// method routes
	router.Get("/get", handlers.HandleMethod)
	router.Post("/post", handlers.HandleMethod)
	router.Put("/put", handlers.HandleMethod)
	router.Patch("/patch", handlers.HandleMethod)
	router.Delete("/delete", handlers.HandleMethod)
	return router
}
