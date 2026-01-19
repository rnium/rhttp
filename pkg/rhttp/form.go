package rhttp

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/url"

	"github.com/rnium/rhttp/internal/build"
)



type File struct {
	Filename    string
	ContentType string
	Data        []byte
}

func (f *File) ToBase64Data() string {
	return build.ToBase64Data(f.ContentType, f.Data)
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

func GetFormData(contentType string, body []byte) (*FormData, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}
	formData := newFormData()
	switch mediaType {
	case "application/x-www-form-urlencoded":
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		for name, value := range values {
			formData.Fields[name] = value[0]
		}
	case "multipart/form-data":
		boundary, ok := params["boundary"]
		if !ok || boundary == "" {
			return nil, ErrNoFormData
		}
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
	default:
		return nil, ErrNoFormData

	}
	return formData, nil
}
