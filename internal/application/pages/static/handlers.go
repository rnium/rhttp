package static

import (
	"bytes"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func fileResponse(path string) *rhttp.Response {
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
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return rhttp.ResponseJSON(400, err.Error())
	}

	reader := bytes.NewReader(buffer)
	res := rhttp.NewChunkedResponse(200, reader)
	_ = res.SetHeader("Content-Type", mimetype)
	return res
}

func webStatic(r *rhttp.Request) *rhttp.Response {
	filename, _ := r.Param("*")
	cwd, _ := os.Getwd()
	fullpath := filepath.Join(cwd, "web/static", filename)

	return fileResponse(fullpath)
}
