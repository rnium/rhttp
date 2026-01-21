package application

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func Start(port uint16) {
	router := getRoutes()
	server := rhttp.Serve(port, router)
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
