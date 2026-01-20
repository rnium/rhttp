package rhttp

import (
	"fmt"
	"strings"
)

func errorResponse(statusCode int, message string) *Response {
	
	res := NewResponse(
		statusCode,
		[]byte(message),
	)
	_ = res.SetHeader("content-type", "text/html")
	return res
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
