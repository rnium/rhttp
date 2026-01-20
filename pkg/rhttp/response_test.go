package rhttp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	res := NewResponse(StatusOK, []byte("foobar"))
	_ = res.SetHeader("content-type", "text/html")
	_ = res.SetHeader("cache-control", "public, max-age=31536000, immutable")

	err := res.SetHeader("Content-length", "2026")
	assert.True(t, errors.Is(err, ErrNonEditableHeader))

	assert.Equal(t, 5, res.headers.Count())
	ctype, _ := res.headers.Get("content-type")
	assert.Equal(t, "text/html", ctype)

	res = NewResponse(StatusOK, []byte("barbaz"))
	assert.Equal(t, 4, res.headers.Count())
}
