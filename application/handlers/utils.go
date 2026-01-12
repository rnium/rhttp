package handlers

import (
	"crypto/tls"
	"net"
	"strings"

	"github.com/rnium/rhttp/internal/http/request"
)

func getClientIP(r *request.Request) string {
	if xff, _ := r.Headers.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}

	if xrip, _ := r.Headers.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr().String())
	return host
}

func detectScheme(req *request.Request) string {
	// check headers
	if proto, ok := req.Headers.Get("X-Forwarded-Proto"); ok && proto != "" {
		return proto
	}
	if scheme, ok := req.Headers.Get("X-Forwarded-Scheme"); ok && scheme != "" {
		return scheme
	}

	// check if connection has tls
	if _, ok := req.Conn().(*tls.Conn); ok {
		return "https"
	}

	return "http"
}

func detectHost(req *request.Request) string {
	if host, ok := req.Headers.Get("Host"); ok && host != "" {
		return host
	}

	// Fallback to local address
	if addr, ok := req.LocalAddr().(*net.TCPAddr); ok {
		return addr.IP.String()
	}

	return "localhost"
}

func normalizeHost(scheme, host string) string {
	if strings.Contains(host, ":") {
		h, p, err := net.SplitHostPort(host)
		if err == nil {
			if (scheme == "http" && p == "80") ||
				(scheme == "https" && p == "443") {
				return h
			}
		}
	}
	return host
}

func buildFullURL(req *request.Request) string {
	scheme := detectScheme(req)
	host := detectHost(req)
	host = normalizeHost(scheme, host)

	return scheme + "://" + host + req.RequestLine.Target
}
