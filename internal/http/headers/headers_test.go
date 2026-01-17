package headers

import (
	"fmt"
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
	h := NewHeaders()
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
	h := NewHeaders()
	headers_test_data := [][2]string{
		{"content-type", "application/json"},
		{"server", "rhttp"},
		{"content-length", "512"}, 
	}
	formatter := func(name, val string) string {
		return fmt.Sprintf("%s: %s\r\n", name, val)
	}
	var expected_output string
	for _, d := range headers_test_data {
		name, val := d[0], d[1]
		_ = h.Set(name, val)
		expected_output += formatter(name, val)
	}
	actual_op := ""
	h.ForEach(func(name, value string) {
		actual_op += formatter(name, value)
	})
	assert.Equal(t, expected_output, actual_op)

}