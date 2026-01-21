package response

import "github.com/rnium/rhttp/pkg/rhttp"

func hasHeader(h *rhttp.Headers, name string) bool {
	_, exists := h.Get(name)
	return exists
}
