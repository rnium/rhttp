package main

import (
	"mime"
	"os"
	"strconv"

	"github.com/rnium/rhttp/internal/application"
)

var PORT uint16 = 8000

func init() {
	portStr := os.Getenv("GO_PORT")	
	if p, err := strconv.Atoi(portStr); err == nil {
		PORT = uint16(p)
	}
	_ = mime.AddExtensionType(".yaml", "application/yaml")
}

func main() {	
	application.Start(PORT)
}
