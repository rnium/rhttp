package methods

import (
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
	"github.com/rnium/rhttp/internal/respond"
)

func handleMethod(r *request.Request) *response.Response {
	var payload any
	if isReadMethod(r.RequestLine.Method) {
		payload = buildReadData(r)
	} else {
		payload = buildWriteData(r)
	}
	return respond.JSON(200, payload)
}
