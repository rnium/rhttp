package methods

import (
	"github.com/rnium/rhttp/pkg/rhttp"
)

func handleMethod(r *rhttp.Request) *rhttp.Response {
	var payload any
	if isReadMethod(r.RequestLine.Method) {
		payload = buildReadData(r)
	} else {
		payload = buildWriteData(r)
	}
	return rhttp.ResponseJSON(200, payload)
}
