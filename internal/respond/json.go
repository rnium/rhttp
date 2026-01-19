package respond

import (
	"encoding/json"
	"errors"

	"github.com/rnium/rhttp/pkg/rhttp"
)

func JSON(status int, payload any) *rhttp.Response {
	data, err := json.Marshal(payload)
	if err != nil {
		return rhttp.Response500(
			errors.New("error transforming data to json"),
		)
	}

	headers := rhttp.NewHeaders()
	_ = headers.Set("Content-Type", "application/json")
	return rhttp.NewResponse(status, data, headers)
}
