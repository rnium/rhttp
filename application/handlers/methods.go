package handlers

import (
	"encoding/json"
	"errors"

	"github.com/rnium/rhttp/internal/http/form"
	"github.com/rnium/rhttp/internal/http/headers"
	"github.com/rnium/rhttp/internal/http/request"
	"github.com/rnium/rhttp/internal/http/response"
	"github.com/rnium/rhttp/internal/utils"
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

func newReadData() *ReadResponseData {
	return &ReadResponseData{
		Args:    make(stringMap),
		Headers: make(stringMap),
	}
}

func newWriteData() *WriteResponseData {
	return &WriteResponseData{
		Args:    make(stringMap),
		Files:   make(stringMap),
		Form:    make(stringMap),
		Headers: make(stringMap),
	}
}

func getReadData(req *request.Request) *ReadResponseData {
	rd := newReadData()
	req.Headers.ForEach(func(name, value string) {
		rd.Headers[name] = value
	})
	req.QParamForEach(func(name, value string) {
		rd.Args[name] = value
	})

	rd.Origin = getClientIP(req)
	rd.Url = buildFullURL(req)
	return rd
}

func getWriteData(req *request.Request) *WriteResponseData {
	wd := newWriteData()
	req.Headers.ForEach(func(name, value string) {
		wd.Headers[name] = value
	})
	req.QParamForEach(func(name, value string) {
		wd.Args[name] = value
	})
	contentType, _ := req.Headers.Get("content-type")
	formdata, _ := form.GetFormData(req)
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
		wd.Data = utils.ToBase64Data(contentType, req.Body)
	}
	wd.Origin = getClientIP(req)
	wd.Url = buildFullURL(req)
	return wd
}

func HandleMethod(r *request.Request) *response.Response {
	var data []byte
	var err error
	if r.RequestLine.Method == "GET" {
		rd := getReadData(r)
		data, err = json.Marshal(rd)
	} else {
		wd := getWriteData(r)
		data, err = json.Marshal(wd)
	}	
	if err != nil {
		return response.Response500(errors.New("Cannot convert data to json"))
	}
	headers := headers.NewHeaders()
	_ = headers.Set("Content-type", "application/json")
	return response.NewResponse(200, data, headers)
}
