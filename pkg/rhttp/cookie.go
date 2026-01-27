package rhttp

import (
	"fmt"
	"strings"
)

type Cookie struct {
	Name     string
	Value    string
	Path     string
	HttpOnly bool
	Secure   bool
	MaxAge   int
}

func formatField(name string, value any) string {
	if _, isBool := value.(bool); isBool {
		return name
	} else {
		return fmt.Sprintf("%s=%s", name, value)
	}
}

func (c *Cookie) String() string {
	parts := []string{
		formatField(c.Name, c.Value),
	}
	if c.Path == "" {
		parts = append(parts, formatField("Path", "/"))
	} else {
		parts = append(parts, formatField("Path", c.Path))
	}

	if c.HttpOnly {
		parts = append(parts, formatField("HttpOnly", c.HttpOnly))
	}

	if c.Secure {
		parts = append(parts, formatField("Secure", c.Secure))
	}

	if c.MaxAge != 0 {
		parts = append(parts, formatField("Max-Age", c.MaxAge))
	}

	return strings.Join(parts, "; ")
}

func (r *Request) Cookies() map[string]string {
	hVal, _ := r.Headers.Get("Cookie")
	cookieMap := make(map[string]string)
	for cRaw := range strings.SplitSeq(hVal, ";") {
		cParts := strings.Split(cRaw, "=")
		if len(cParts) != 2 {
			continue
		}
		cookieMap[strings.TrimSpace(cParts[0])] = strings.TrimSpace(cParts[1])
	}
	return cookieMap
}
