package static

import (
	"os"
	"path/filepath"

	"github.com/rnium/rhttp/internal/respond"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func webStatic(r *rhttp.Request) *rhttp.Response {
	filename, _ := r.Param("*")
	cwd, _ := os.Getwd()
	fullpath := filepath.Join(cwd, "web/static", filename)

	return respond.FileResponse(fullpath)
}
