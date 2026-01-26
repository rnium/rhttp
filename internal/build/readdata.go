package build

import (
	"github.com/rnium/rhttp/internal/inspect"
	"github.com/rnium/rhttp/pkg/rhttp"
)

type ReadResponseData struct {
	Args    Args    `json:"args"`
	Headers Headers `json:"headers"`
	Origin  string  `json:"origin"`
	Url     string  `json:"url"`
}

func BuildReadData(req *rhttp.Request) *ReadResponseData {
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
