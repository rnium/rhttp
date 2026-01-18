package requestinspection

import (
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
	"github.com/rnium/rhttp/internal/inspect"
	"github.com/rnium/rhttp/internal/respond"
)

func viewHeaders(r *request.Request) *response.Response {
	headersData := buildHeadersData(r.Headers)
	return respond.JSON(200, headersData)
}

func viewIp(r *request.Request) *response.Response {
	data := make(map[string]string)
	data["origin"] = inspect.ClientIP(r)

	return respond.JSON(200, data)
}

func viewUserAgent(r *request.Request) *response.Response {
	data := make(map[string]any)
	userAgent, exists := r.Headers.Get("user-agent")
	if exists {
		data["user-agent"] = userAgent
	} else {
		data["user-agent"] = nil
	}
	return respond.JSON(200, data)
}
