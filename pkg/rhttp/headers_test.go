package rhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	assert.NoError(t, validateToken("content-length"))
	assert.ErrorIs(t, validateToken("content-(length)"), ErrInvalidToken)
	assert.ErrorIs(t, validateToken("contentlength:"), ErrInvalidToken)
	assert.ErrorIs(t, validateToken(""), ErrEmptyToken)
	assert.ErrorIs(t, validateToken("content length"), ErrInvalidToken)
}

func TestHeaders(t *testing.T) {
	h := newHeaders()
	err := h.Set("server", "nginx")
	require.NoError(t, err)
	val, ok := h.Get("server")
	assert.True(t, ok)
	assert.Equal(t, val, "nginx")
	// with invalid token
	err = h.Set("content length", "123")
	require.ErrorIs(t, err, ErrInvalidToken)
	// multi value
	_ = h.Set("foo", "bar")
	_ = h.Set("foo", "baz")
	val, _ = h.Get("foo")
	assert.Equal(t, val, "bar, baz")
	// replace
	_, new := h.Replace("foo", "abc")
	assert.False(t, new)
	val, exists := h.Get("foo")
	assert.True(t, exists)
	assert.Equal(t, val, "abc")
	err, new = h.Replace("nonexistant", "somevalue")
	assert.NoError(t, err)
	assert.True(t, new)
	val, _ = h.Get("nonexistant")
	assert.Equal(t, val, "somevalue")
	// remove
	h.Remove("foo")
	_, ok = h.Get("foo")
	assert.False(t, ok)
}

func TestHeadersForEach(t *testing.T) {
	h := newHeaders()
	headers_test_data := [][2]string{
		{"content-type", "application/json"},
		{"server", "rhttp"},
		{"content-length", "512"},
	}

	for _, d := range headers_test_data {
		name, val := d[0], d[1]
		_ = h.Set(name, val)
	}

	for _, d := range headers_test_data {
		name, val_exp := d[0], d[1]
		val_actual, _ := h.Get(name)
		assert.Equal(t, val_exp, val_actual)
	}

	var header_names []string
	h.ForEach(func(name, value string) {
		header_names = append(header_names, name)
	})

	assert.Equal(t, 3, len(header_names))

}
