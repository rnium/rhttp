package respond

import (
	"encoding/json"
	"errors"

	"github.com/rnium/rhttp/internal/http/headers"
	"github.com/rnium/rhttp/internal/http/response"
)

func JSON(status int, payload any) *response.Response {
	data, err := json.Marshal(payload)
	if err != nil {
		return response.Response500(
			errors.New("error transforming data to json"),
		)
	}

	headers := headers.NewHeaders()
	_ = headers.Set("Content-Type", "application/json")
	return response.NewResponse(status, data, headers)
}
