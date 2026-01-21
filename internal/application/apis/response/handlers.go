package response

import (
	"fmt"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func cache(r *rhttp.Request) *rhttp.Response {
	statusCode := 200
	var payload any
	if hasHeader(r.Headers, "if-modified-since") || hasHeader(r.Headers, "if-none-match") {
		statusCode = 304
		payload = nil
	} else {
		payload = buildReadData(r)
	}
	if payload == nil {
		return rhttp.NewResponse(statusCode, nil)
	}
	return rhttp.ResponseJSON(statusCode, payload)
}

func setCacheCtrl(r *rhttp.Request) *rhttp.Response {
	value, _ := r.Param("value")
	if !isInteger(value) {
		payload := map[string]string{
			"message": "value should be an integer",
		}
		return rhttp.ResponseJSON(400, payload)
	}
	payload := buildReadData(r)
	res := rhttp.ResponseJSON(200, payload)
	_ = res.SetHeader("cache-control", fmt.Sprintf("public, max-age=%s", value))
	return res
}
