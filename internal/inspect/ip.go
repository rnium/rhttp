package inspect

import (
	"net"
	"strings"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func ClientIP(r *rhttp.Request) string {
	if xff, _ := r.Headers.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}

	if xrip, _ := r.Headers.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr().String())
	return host
}
