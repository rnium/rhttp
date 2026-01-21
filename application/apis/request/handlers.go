package request

import (
	"github.com/rnium/rhttp/internal/inspect"
	"github.com/rnium/rhttp/pkg/rhttp"
)

func viewHeaders(r *rhttp.Request) *rhttp.Response {
	headersData := buildHeadersData(r.Headers)
	return rhttp.ResponseJSON(200, headersData)
}

func viewIp(r *rhttp.Request) *rhttp.Response {
	data := make(map[string]string)
	data["origin"] = inspect.ClientIP(r)

	return rhttp.ResponseJSON(200, data)
}

func viewUserAgent(r *rhttp.Request) *rhttp.Response {
	data := make(map[string]any)
	userAgent, exists := r.Headers.Get("user-agent")
	if exists {
		data["user-agent"] = userAgent
	} else {
		data["user-agent"] = nil
	}
	return rhttp.ResponseJSON(200, data)
}
