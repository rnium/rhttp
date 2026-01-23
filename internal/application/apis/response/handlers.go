package response

import (
	"fmt"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func statusResponse(statusCode int, payload any, etag string) *rhttp.Response {
	var res *rhttp.Response
	if payload == nil {
		res = rhttp.NewResponse(statusCode, nil)
	} else {
		res = rhttp.ResponseJSON(statusCode, payload)
	}
	if etag != "" {
		_ = res.SetHeader("etag", etag)
	}
	return res
}

func cache(r *rhttp.Request) *rhttp.Response {
	statusCode := 200
	var payload any
	if hasHeader(r.Headers, "if-modified-since") || hasHeader(r.Headers, "if-none-match") {
		statusCode = 304
		payload = nil
	} else {
		payload = buildReadData(r)
	}
	return statusResponse(statusCode, payload, "")
}

func setCacheCtrl(r *rhttp.Request) *rhttp.Response {
	value, _ := r.Param("value")
	if !isInteger(value) {
		payload := map[string]string{
			"message": "value should be an integer",
		}
		return rhttp.ResponseJSON(400, payload)
	}
	payload := buildReadData(r)
	res := rhttp.ResponseJSON(200, payload)
	_ = res.SetHeader("cache-control", fmt.Sprintf("public, max-age=%s", value))
	return res
}

func etagHandler(r *rhttp.Request) *rhttp.Response {
	etag, _ := r.Param("etag")
	if etag == "" {
		etag = "{etag}"
	}

	statusCode := 200 
	var payload any = buildReadData(r)

	ifNoneMatch, ok := r.Headers.Get("if-none-match")
	if ok {
		if ifNoneMatch == etag {
			statusCode = 304
			payload = nil
		}
		return statusResponse(statusCode, payload, etag)
	}
	ifMatch, ok := r.Headers.Get("if-match")
	if ok {
		if ifMatch != etag {
			statusCode = 412
			payload = nil
		}
	}
	return statusResponse(statusCode, payload, etag)
}
