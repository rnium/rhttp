package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rnium/rhttp/internal/server"
)




const PORT uint16 = 8980


func main() {
	server := server.Serve(PORT)
	defer server.Close()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	s := <-sigChan
	fmt.Printf("{%v}\n", s)
	fmt.Println("Shutting down gracefully")
}