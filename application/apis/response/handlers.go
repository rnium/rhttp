package response

import (
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
	"github.com/rnium/rhttp/internal/respond"
)

func cache(r *request.Request) *response.Response {
	statusCode := 200
	var payload any
	if hasHeader(r.Headers, "if-modified-since") || hasHeader(r.Headers, "if-none-match") {
		statusCode = 304
		payload = nil
	} else {
		payload = buildReadData(r)
	}
	if payload == nil {
		return response.NewResponse(statusCode, nil, nil)
	}
	return respond.JSON(statusCode, payload)
}
