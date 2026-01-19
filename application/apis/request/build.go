package request

import "github.com/rnium/rhttp/pkg/rhttp"

type HeadersData struct {
	Headers map[string]string `json:"headers"`
}

func buildHeadersData(h *rhttp.Headers) *HeadersData {
	dataMap := make(map[string]string)
	h.ForEach(func(name, value string) {
		dataMap[name] = value
	})
	return &HeadersData{
		Headers: dataMap,
	}
}
