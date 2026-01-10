package router

import (
	"fmt"
	"net/url"
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


func parseTarget(target string) (string, request.Params) {
	parts := strings.SplitN(target, "?", 2)
	if len(parts) == 1 {
		return parts[0], nil
	}

	queryParams := request.NewParams()

	for pair := range strings.SplitSeq(parts[1], "&") {
		if pair == "" {
			continue
		}

		key, val, found := strings.Cut(pair, "=")
		if !found {
			continue
		}

		k, err1 := url.QueryUnescape(key)
		v, err2 := url.QueryUnescape(val)
		if err1 != nil || err2 != nil {
			continue
		}

		queryParams[k] = v
	}

	return parts[0], queryParams
}
