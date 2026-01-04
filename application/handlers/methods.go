package handlers

import (
	"encoding/base64"
	"fmt"

	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/response"
)

func HandlePost(r *request.Request) *response.Response {
	ctype, _ := r.Headers.Get("content-type")
	res_str := ""
	if ctype == "image/png" {
		res_str = fmt.Sprintf("data: %s", base64.StdEncoding.EncodeToString(r.Body))
	}
	return response.NewResponse(200, []byte(res_str), nil)
}
