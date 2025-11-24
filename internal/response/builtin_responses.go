package response

import "fmt"

var Response500 = func(err error) *Response {
	msg := fmt.Sprint(err)
	return NewResponse(
		StatusInternalServerError,
		[]byte(msg),
		nil,
	)
}
