package response

import (
	"github.com/rnium/rhttp/internal/respond"
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
	return respond.JSON(statusCode, payload)
}


// func setCacheCtrl(r *rhttp.Request) *rhttp.Response {
// 	val, ok := r.Param("value")
	
// }