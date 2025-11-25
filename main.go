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
	router.Get("/health", HealthCheck)
	router.Get("/servererror", MyBad)
	router.Get("/param/:pk/info", ParamChecker)
	router.Get("/method-check", MethodChecker)
	router.Post("/method-check", MethodChecker)
	router.Put("/method-check", MethodChecker)
	router.Patch("/method-check", MethodChecker)
	router.Delete("/method-check", MethodChecker)
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