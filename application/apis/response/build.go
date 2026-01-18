package response

import (
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/inspect"
)

type Args map[string]string
type Headers map[string]string
type Form map[string]string
type Files map[string]string

type ReadResponseData struct {
	Args    Args    `json:"args"`
	Headers Headers `json:"headers"`
	Origin  string  `json:"origin"`
	Url     string  `json:"url"`
}

func buildReadData(req *request.Request) *ReadResponseData {
	rd := &ReadResponseData{
		Args:    make(Args),
		Headers: make(Headers),
	}
	req.Headers.ForEach(func(name, value string) {
		rd.Headers[name] = value
	})
	req.QParamForEach(func(name, value string) {
		rd.Args[name] = value
	})

	rd.Origin = inspect.ClientIP(req)
	rd.Url = inspect.FullURL(req)
	return rd
}
