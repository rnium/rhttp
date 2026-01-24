package dynamicdata

import (
	"encoding/base64"
	"fmt"

	"github.com/rnium/rhttp/pkg/rhttp"
)

var DEFAULT_BASE64_DATA = "ckh0dHBiaW4gaXMgYXdlc29tZQ=="

func decodeBase64Handler(r *rhttp.Request) *rhttp.Response {
	value, _ := r.Param("value")
	payload, err := base64.StdEncoding.DecodeString(value)
	statusCode := rhttp.StatusOK
	if err != nil {
		payload = fmt.Appendf(nil, "Incorrect Base64 data try: %s", DEFAULT_BASE64_DATA)
		statusCode = rhttp.StatusBadRequest
	}
	res := rhttp.NewResponse(statusCode, payload)
	_ = res.SetHeader("Content-type", "text/plain")
	return res
}

func uuidGenHandler(_ *rhttp.Request) *rhttp.Response {
	uid := buildUUID()
	data := map[string]string {
		"uuid": uid,
	}
	return rhttp.ResponseJSON(200, data)
}