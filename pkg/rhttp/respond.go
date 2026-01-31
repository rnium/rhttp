// Formerly known as builtin responses
// A Collection of common responses
package rhttp

import (
	"encoding/json"
	"strings"
)

func ResponseJSON(status int, payload any) *Response {
	data, err := json.Marshal(payload)
	if err != nil {
		status = 500
		data, _ = json.Marshal(map[string]string{
			"message": "error transforming data to json",
		})
	}

	res := NewResponse(status, data)
	_ = res.SetHeader("Content-Type", "application/json")
	return res
}

func ErrorResponseJSON(status int, messages ...string) *Response {
	var msg string
	if len(messages) > 0 {
		msg = strings.Join(messages, ", ")
	} else {
		msg = statusMessage[status]
	}
	return ResponseJSON(
		status,
		map[string]string{
			"error": msg,
		},
	)
}

func response500(err error) *Response {
	return ErrorResponseJSON(500, err.Error())
}

func Redirect(to string, parmanent bool) *Response {
	statusCode := 302
	if parmanent {
		statusCode = 301
	}
	res := ResponseJSON(statusCode, "")
	_ = res.SetHeader("Location", to)
	return res
}
