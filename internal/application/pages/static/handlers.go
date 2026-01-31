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

	res := respond.FileResponse(fullpath)
	_ = res.SetHeader("cache-control", "public, max-age=315360000, immutable")
	return res
}
