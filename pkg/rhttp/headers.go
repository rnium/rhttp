package rhttp

import (
	"strings"
)

type Header struct {
	Name  string
	Value string
}

func newHeader(name, value string) *Header {
	return &Header{
		Name:  name,
		Value: value,
	}
}

type Headers struct {
	headers map[string]*Header
}

func validateToken(token string) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	for _, c := range token {
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case strings.ContainsRune("!#$%&'*+-.^_`|~", c): // RFC 9110 #5.6.2
		default:
			return ErrInvalidToken
		}
	}
	return nil
}

func (h *Headers) Set(name, value string) error {
	err := validateToken(name)
	if err != nil {
		return err
	}
	name_lower := strings.ToLower(name)
	if prevHeader, exists := h.headers[name_lower]; exists {
		value = strings.Join([]string{prevHeader.Value, value}, ", ") // RFC 9110 #5.3
		prevHeader.Value = value
	} else {
		h.headers[name_lower] = newHeader(name, value)
	}
	return nil
}

func (h *Headers) Get(name string) (string, bool) {
	val, ok := h.headers[strings.ToLower(name)]
	if ok {
		return val.Value, ok
	}
	return "", ok
}

func (h *Headers) Replace(name, newValue string) (err error, new bool) {
	name_lower := strings.ToLower(name)
	if _, exists := h.headers[name_lower]; !exists {
		err = h.Set(name_lower, newValue)
		return err, true
	}
	h.headers[name_lower] = newHeader(name, newValue)
	return nil, false
}

func (h *Headers) Remove(name string) {
	delete(h.headers, strings.ToLower(name))
}

func (h *Headers) ForEach(f func(name, value string)) {
	for _, header := range h.headers {
		f(header.Name, header.Value)
	}
}

func (h *Headers) Count() int {
	return len(h.headers)
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]*Header),
	}
}

func GetDefaultResponseHeaders() *Headers {
	headers := NewHeaders()
	defaults := map[string]string{
		"content-type": "text/plain",
		"server":       "rhttp",
		"connection":   "closed",
	}
	for name, value := range defaults {
		_ = headers.Set(name, value)
	}
	return headers
}
