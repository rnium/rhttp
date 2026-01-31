package rhttp

const (
	StatusOK                          = 200 // RFC 9110, 15.3.1
	StatusCreated                     = 201 // RFC 9110, 15.3.2
	StatusAccepted                    = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInformation = 203 // RFC 9110, 15.3.4
	StatusNoContent                   = 204 // RFC 9110, 15.3.5

	StatusMultipleChoices  = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently = 301 // RFC 9110, 15.4.2
	StatusFound            = 302 // RFC 9110, 15.4.3

	StatusBadRequest       = 400 // RFC 9110, 15.5.1
	StatusUnauthorized     = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired  = 402 // RFC 9110, 15.5.3
	StatusForbidden        = 403 // RFC 9110, 15.5.4
	StatusNotFound         = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable    = 406 // RFC 9110, 15.5.7

	StatusInternalServerError     = 500 // RFC 9110, 15.6.1
	StatusNotImplemented          = 501 // RFC 9110, 15.6.2
	StatusBadGateway              = 502 // RFC 9110, 15.6.3
	StatusServiceUnavailable      = 503 // RFC 9110, 15.6.4
	StatusGatewayTimeout          = 504 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported = 505 // RFC 9110, 15.6.6
)

var statusMessage = map[int]string{
	StatusOK:                          "OK",
	StatusCreated:                     "Created",
	StatusAccepted:                    "Accepted",
	StatusNonAuthoritativeInformation: "Non-Authoritative Information",
	StatusNoContent:                   "No Content",

	StatusMultipleChoices:  "Multiple Choices",
	StatusMovedPermanently: "Moved Permanently",
	StatusFound:            "Found",

	StatusBadRequest:       "Bad Request",
	StatusUnauthorized:     "Unauthorized",
	StatusPaymentRequired:  "Payment Required",
	StatusForbidden:        "Forbidden",
	StatusNotFound:         "Not Found",
	StatusMethodNotAllowed: "Method Not Allowed",
	StatusNotAcceptable:    "Not Acceptable",

	StatusInternalServerError:     "Internal Server Error",
	StatusNotImplemented:          "Not Implemented",
	StatusBadGateway:              "Bad Gateway",
	StatusServiceUnavailable:      "Service Unavailable",
	StatusGatewayTimeout:          "Gateway Timeout",
	StatusHTTPVersionNotSupported: "HTTP Version Not Supported",
}
