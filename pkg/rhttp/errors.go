package rhttp

import (
	"errors"
)

// Form Errors
var ErrNoFormData = errors.New("request has no form data")
var ErrReadingFormPart = errors.New("error reading form part")
var ErrReadingFileData = errors.New("error while reading filedata")

// Headers
var ErrInvalidToken = errors.New("token contains invalid characters")
var ErrEmptyToken = errors.New("token is empty")

// Request
var ErrMalformedRequestLine = errors.New("malformed request line")
var ErrMalformedFieldLine = errors.New("malformed field line")
var ErrMalformedFieldValue = errors.New("malformed field value")

// Response
var ErrResponseClosed = errors.New("response closed")
var ErrWriteFailed = errors.New("error while writing to connection")
