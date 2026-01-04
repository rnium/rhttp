package main

import (
	"github.com/rnium/rhttp/application"
)

const PORT uint16 = 8000

func main() {
	application.Start(PORT)
}
