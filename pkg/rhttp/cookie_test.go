package rhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieString(t *testing.T) {
	cookie := &Cookie{
		Name: "session_id",
		Value: "abc",
	}
	assert.Equal(t, "session_id=abc; Path=/", cookie.String())
	cookie = &Cookie{
		Name: "Mode",
		Value: "dark",
		Path: "/dashboard",
		Secure: true,
		HttpOnly: true,
	}
	assert.Equal(t, "Mode=dark; Path=/dashboard; HttpOnly; Secure", cookie.String())
}