package methods

import (
	"encoding/json"

	"github.com/rnium/rhttp/internal/build"
	"github.com/rnium/rhttp/internal/inspect"
	"github.com/rnium/rhttp/pkg/rhttp"
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

type WriteResponseData struct {
	Args    Args           `json:"args"`
	Data    string         `json:"data"`
	Files   Files          `json:"files"`
	Form    Form           `json:"form"`
	Headers Headers        `json:"headers"`
	Json    map[string]any `json:"json"`
	Origin  string         `json:"origin"`
	Url     string         `json:"url"`
}

func buildReadData(req *rhttp.Request) *ReadResponseData {
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

func buildWriteData(req *rhttp.Request) *WriteResponseData {
	wd := &WriteResponseData{
		Args:    make(Args),
		Files:   make(Files),
		Form:    make(Form),
		Headers: make(Headers),
	}
	req.Headers.ForEach(func(name, value string) {
		wd.Headers[name] = value
	})
	req.QParamForEach(func(name, value string) {
		wd.Args[name] = value
	})
	contentType, _ := req.Headers.Get("content-type")
	formdata, _ := req.FormData()
	if formdata != nil {
		wd.Form = formdata.Fields
		for field, file := range formdata.Files {
			wd.Files[field] = file.ToBase64Data()
		}
	} else if contentType == "application/json" {
		wd.Data = string(req.Body)
		jsonData := make(map[string]any)
		_ = json.Unmarshal(req.Body, &jsonData)
		wd.Json = jsonData
	} else if len(req.Body) > 0 {
		wd.Data = build.ToBase64Data(contentType, req.Body)
	}
	wd.Origin = inspect.ClientIP(req)
	wd.Url = inspect.FullURL(req)
	return wd
}
