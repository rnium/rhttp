package form

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"

	"github.com/rnium/rhttp/internal/http/request"
)

var ErrNoFormData = errors.New("Request has no form data")
var ErrReadingFormPart = errors.New("Error reading form part")
var ErrReadingFileData = errors.New("Error while reading filedata")

type File struct {
	Filename    string
	ContentType string
	Data        []byte
}

func (f *File) ToBase64Data() string {
	return fmt.Sprintf(
		"data:%s;base64,%s",
		f.ContentType,
		base64.StdEncoding.EncodeToString(f.Data),
	)
}

type FormData struct {
	Fields map[string]string
	Files  map[string]*File
}

// func (fd *FormData) AddField(name, value string) {
// 	fd.Fields[name] = value
// }

func newFormData() *FormData {
	return &FormData{
		Fields: make(map[string]string),
		Files:  make(map[string]*File),
	}
}

func getFormData(contentType string, body []byte) (*FormData, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "multipart/form-data" {
		return nil, ErrNoFormData
	}
	boundary, ok := params["boundary"]
	if !ok || boundary == "" {
		return nil, ErrNoFormData
	}

	formData := newFormData()
	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Join(ErrReadingFormPart, err)
		}
		if part.FileName() == "" {
			value, _ := io.ReadAll(part)
			formData.Fields[part.FormName()] = string(value)
			continue
		}
		data, err := io.ReadAll(part)
		if err != nil {
			return nil, ErrReadingFileData
		}
		formData.Files[part.FormName()] = &File{
			Filename:    part.FileName(),
			ContentType: part.Header.Get("Content-Type"),
			Data:        data,
		}
	}
	return formData, nil
}

func GetFormData(r *request.Request) (*FormData, error) {
	if r == nil {
		return nil, ErrNoFormData
	}
	ctype, _ := r.Headers.Get("content-type")
	return getFormData(ctype, r.Body)
}
