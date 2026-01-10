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

type WriteData struct {
	Args    stringMap      `json:"args"`
	Data    string         `json:"data"`
	Files   stringMap      `json:"files"`
	Form    stringMap      `json:"form"`
	Headers stringMap      `json:"headers"`
	Json    map[string]any `json:"json"`
	Origin  string         `json:"origin"`
	Url     string         `json:"url"`
}

func newWriteData() *WriteData {
	return &WriteData{
		Args:    make(stringMap),
		Files:   make(stringMap),
		Form:    make(stringMap),
		Headers: make(stringMap),
	}
}

func getWriteData(req *request.Request) *WriteData {
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
		jsonData := make(map[string]any)
		_ = json.Unmarshal(req.Body, &jsonData)
		wd.Json = jsonData
	} else {
		wd.Data = utils.ToBase64Data(contentType, req.Body)
	}
	return wd
}

func HandlePost(r *request.Request) *response.Response {
	wd := getWriteData(r)
	data, err := json.Marshal(wd)
	if err != nil {
		return response.Response500(errors.New("Cannot convert data to json"))
	}
	headers := headers.NewHeaders()
	_ = headers.Set("Content-type", "application/json")
	return response.NewResponse(200, data, headers)
}
