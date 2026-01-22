package main

import (
	"mime"

	"github.com/rnium/rhttp/internal/application"
)

const PORT uint16 = 8000

func init() {
	_ = mime.AddExtensionType(".yaml", "application/yaml")
}

func main() {
	application.Start(PORT)
}
