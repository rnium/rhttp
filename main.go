package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rnium/rhttp/internal/router"
	"github.com/rnium/rhttp/internal/server"
)

const PORT uint16 = 8000

func getRouter() *router.Router {
	router := router.NewRouter()
	router.Get("/", SwaggerUI)
	router.Get("/ping", Ping)
	router.Get("/swagger/openapi.yaml", OpenAPISpec)
	router.Get("/health", HealthCheck)
	router.Get("/httpbin/stream/:n", HttpBinStream)
	return router
}

func main() {
	router := getRouter()
	server := server.Serve(PORT, router)
	defer func() {
		fmt.Println("Shutting down gracefully")
		err := server.Close()
		if err != nil {
			fmt.Printf("Error while shutting down server %v\n", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
