package methods

import (
	"github.com/rnium/rhttp/internal/build"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func handleMethod(r *rhttp.Request) *rhttp.Response {
	var payload any
	if isReadMethod(r.RequestLine.Method) {
		payload = build.BuildReadData(r)
	} else {
		payload = build.BuildWriteData(r)
	}
	return rhttp.ResponseJSON(200, payload)
}
