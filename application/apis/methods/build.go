package methods

import (
	"encoding/json"

	"github.com/rnium/rhttp/internal/build"
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/inspect"
)

type stringMap map[string]string

type ReadResponseData struct {
	Args    stringMap `json:"args"`
	Headers stringMap `json:"headers"`
	Origin  string    `json:"origin"`
	Url     string    `json:"url"`
}

type WriteResponseData struct {
	Args    stringMap      `json:"args"`
	Data    string         `json:"data"`
	Files   stringMap      `json:"files"`
	Form    stringMap      `json:"form"`
	Headers stringMap      `json:"headers"`
	Json    map[string]any `json:"json"`
	Origin  string         `json:"origin"`
	Url     string         `json:"url"`
}

func newReadResponseData() *ReadResponseData {
	return &ReadResponseData{
		Args:    make(stringMap),
		Headers: make(stringMap),
	}
}

func newWriteResponseData() *WriteResponseData {
	return &WriteResponseData{
		Args:    make(stringMap),
		Files:   make(stringMap),
		Form:    make(stringMap),
		Headers: make(stringMap),
	}
}

func buildReadData(req *request.Request) *ReadResponseData {
	rd := newReadResponseData()
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


func buildWriteData(req *request.Request) *WriteResponseData {
	wd := newWriteResponseData()
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

