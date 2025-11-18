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
		case strings.ContainsRune("!#$%&'*+-.^_`|~", c):
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
		value = strings.Join([]string{prev_val, value}, ",")
	}
	h.headers[name_lower] = value
	return nil
}

func (h *Headers) Get(name string) (string, bool) {
	val, ok := h.headers[strings.ToLower(name)]
	return val, ok
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}