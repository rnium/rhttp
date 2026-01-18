package response

import "github.com/rnium/rhttp/internal/http/headers"

func hasHeader(h *headers.Headers, name string) bool {
	_, exists := h.Get(name)
	return exists
}