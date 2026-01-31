package rhttp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCookieString(t *testing.T) {
	cookie := &Cookie{
		Name:  "session_id",
		Value: "abc",
	}
	assert.Equal(t, "session_id=abc; Path=/", cookie.String())
	cookie = &Cookie{
		Name:     "Mode",
		Value:    "dark",
		Path:     "/dashboard",
		Secure:   true,
		HttpOnly: true,
	}
	assert.Equal(t, "Mode=dark; Path=/dashboard; HttpOnly; Secure", cookie.String())
	expiry, _ := time.Parse(
		time.RFC1123,
		"Tue, 01 Jan 2030 12:00:00 GMT",
	)
	cookie = &Cookie{
		Name:    "free",
		Value:   "form",
		Expires: expiry,
	}
	assert.Equal(t, "free=form; Expires=Tue, 01 Jan 2030 12:00:00 GMT; Path=/", cookie.String())
}
