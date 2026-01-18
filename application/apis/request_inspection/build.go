package requestinspection

import "github.com/rnium/rhttp/internal/http/headers"

type HeadersData struct {
	Headers map[string]string `json:"headers"`
}

func buildHeadersData(h *headers.Headers) *HeadersData {
	dataMap := make(map[string]string)
	h.ForEach(func(name, value string) {
		dataMap[name] = value
	})
	return &HeadersData{
		Headers: dataMap,
	}
}
