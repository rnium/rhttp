package methods

import (
	"encoding/json"
	"errors"

	"github.com/rnium/rhttp/internal/http/headers"
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
)

func handleMethod(r *request.Request) *response.Response {
	var data []byte
	var err error
	if r.RequestLine.Method == "GET" {
		rd := buildReadData(r)
		data, err = json.Marshal(rd)
	} else {
		wd := buildWriteData(r)
		data, err = json.Marshal(wd)
	}
	if err != nil {
		return response.Response500(errors.New("Cannot convert data to json"))
	}
	headers := headers.NewHeaders()
	_ = headers.Set("Content-type", "application/json")
	return response.NewResponse(200, data, headers)
}
