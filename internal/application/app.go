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
	server := rhttp.NewServer(router)
	server.Start(port)
	defer func() {
		fmt.Printf("\nShutting down server instance at %v\n", &server)
		err := server.Close()
		if err != nil {
			fmt.Printf("Error while shutting down server %v\n", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
