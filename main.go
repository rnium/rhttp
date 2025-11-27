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
	router.Get("/", Index)
	router.Get("/ping", Ping)
	router.Get("/health", HealthCheck)
	router.Get("/httpbin/stream/:n", HttpBinStream)
	return router
}

func main() {
	router := getRouter()
	server := server.Serve(PORT, router)
	defer server.Close()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	s := <-sigChan
	fmt.Printf("{%v}\n", s)
	fmt.Println("Shutting down gracefully")
}
