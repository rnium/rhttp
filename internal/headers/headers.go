package headers

import (
	"fmt"
	"strings"
)

var ErrInvalidToken = fmt.Errorf("token contains invalid characters")
var ErrEmptyToken = fmt.Errorf("token is empty")


type Headers struct {
	headers map[string]string
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
		case strings.ContainsRune("!#$%&'*+-.^_`|~", c):		// RFC 9110 #5.6.2
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
	if prev_val, exists := h.headers[name_lower]; exists {
		value = strings.Join([]string{prev_val, value}, ", ") // RFC 9110 #5.3
	}
	h.headers[name_lower] = value
	return nil
}

func (h *Headers) Get(name string) (string, bool) {
	val, ok := h.headers[strings.ToLower(name)]
	return val, ok
}

func (h *Headers) Replace(name, newValue string) (err error, new bool) {
	name_lower := strings.ToLower(name)	
	if _, exists := h.headers[name_lower]; !exists {
		err = h.Set(name_lower, newValue)
		return err, true
	}
	h.headers[name_lower] = newValue
	return nil, false
}


func (h *Headers) Remove(name string) {
	delete(h.headers, strings.ToLower(name))
}

func (h *Headers) ForEach(f func (name, value string)) {
	for k, v := range h.headers {
		f(k, v)
	}
}



func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}