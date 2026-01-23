package respond

import (
	"mime"
	"os"
	"path/filepath"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func FileResponse(path string) *rhttp.Response {
	mimetype := mime.TypeByExtension(filepath.Ext(path))

	file, err := os.Open(path)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return rhttp.ResponseJSON(404, "File not found")
		default:
			return rhttp.ResponseJSON(500, err.Error())
		}
	}

	res := rhttp.NewChunkedResponse(200, file)
	_ = res.SetHeader("Content-Type", mimetype)
	return res
}
