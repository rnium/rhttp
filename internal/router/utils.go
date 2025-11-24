package router

import (
	"fmt"

	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/response"
)

var ErrInvalidHttpTarget = fmt.Errorf("invalid http target")


func validateTarget(target string) error {
	// Todo: Implement url pattern validation
	return nil
}

func NewErrorHandler(statusCode int) Handler {
	return func(r *request.Request) *response.Response {
		return response.ErrorResponseHTML(statusCode)
	}
}