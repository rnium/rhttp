package response

import (
	"fmt"
	"strings"

	"github.com/rnium/rhttp/internal/http/headers"
)

func errorResponse(statusCode int, message string) *Response {
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	return NewResponse(
		statusCode,
		[]byte(message),
		headers,
	)
}

var Response500 = func(err error) *Response {
	msg := fmt.Sprintf("<h1>Internal Server Error</h1><h3 style='color: gray'>Details: %s</h3>", err)
	return errorResponse(StatusInternalServerError, msg)
}

var ErrorResponseHTML = func(status int, messages ...string) *Response {
	var msg string 
	if len(messages) > 0 {
		msg = fmt.Sprintf("<h1>%s</h1>", strings.Join(messages, ", "))
	} else {
		msg = fmt.Sprintf("<h1>%s</h1>", statusMessage[status])
	}
	return errorResponse(status, msg)
}
