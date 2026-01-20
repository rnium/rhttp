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

	res := rhttp.NewResponse(status, data)
	_ = res.SetHeader("Content-Type", "application/json")
	return res
}
