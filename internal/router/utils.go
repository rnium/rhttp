package router

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
)

var ErrInvalidHttpTarget = fmt.Errorf("invalid http target")


func validateTarget(target string) error {
	for _, c := range target {
		switch {
		case unicode.IsLetter(c) || unicode.IsDigit(c):
		case strings.ContainsRune("./-_:", c):
		default:
			return ErrInvalidHttpTarget
		}
	}
	return nil
}

func NewErrorHandler(statusCode int) Handler {
	return func(r *request.Request) *response.Response {
		return response.ErrorResponseHTML(statusCode)
	}
}