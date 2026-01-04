package response

import (
	"fmt"
	"testing"

	"github.com/rnium/rhttp/internal/http/headers"
	"github.com/stretchr/testify/assert"
)


func TestNewResponse(t *testing.T) {
	headers := headers.NewHeaders()
	_ = headers.Set("content-type", "text/html")
	_ = headers.Set("cache-control", "public, max-age=31536000, immutable")
	res := NewResponse(StatusOK, []byte("foobar"), headers)
	res.Headers.ForEach(func(name, value string) {
		fmt.Println(name)
	})
	assert.Equal(t, 5, res.Headers.Count())
	ctype, _ := res.Headers.Get("content-type")
	assert.Equal(t, "text/html", ctype)

	res = NewResponse(StatusOK, []byte("barbaz"), nil)
	assert.Equal(t, 4, res.Headers.Count())
}