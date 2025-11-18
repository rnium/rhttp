package headers

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
	h := NewHeaders()
	err := h.Set("server", "nginx")
	require.NoError(t, err)
	val, ok := h.Get("server")
	assert.True(t, ok)
	assert.Equal(t, val, "nginx")
	// with invalid token
	err = h.Set("content length", "123")
	require.Error(t, err)
}