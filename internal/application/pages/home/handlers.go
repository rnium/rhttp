package home

import (
	"github.com/rnium/rhttp/internal/respond"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func healthCheck(r *rhttp.Request) *rhttp.Response {
	return rhttp.ResponseJSON(
		200,
		map[string]any{
			"success": true,
			"message": "Everything working fine",
		},
	)
}

func index(r *rhttp.Request) *rhttp.Response {
	return respond.FileResponse("./web/templates/index.html")
}

func schema(r *rhttp.Request) *rhttp.Response {
	return respond.FileResponse("./web/openapi.yaml")
}
