package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkedReader struct {
	data []byte
	bytesPerRead int
	pos int
}


func (cr *chunkedReader) Read(p []byte) (int, error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	n := min(cr.bytesPerRead, len(cr.data)-cr.pos)	
	n = copy(p, cr.data[cr.pos:cr.pos+n])
	cr.pos += n
	
	return n, nil
}

func newChunkedReader(data []byte) *chunkedReader {
	return &chunkedReader{
		data: data,
		bytesPerRead: 2,
		pos: 0,
	}
}


func TestRequest(t *testing.T) {
	cr := newChunkedReader([]byte("GET hello/ HTTP1.1\r\ncontent-length: 43\r\naccept: text/plain\r\naccept: application/json\r\n\r\n"))
	rq, err := GetRequest(cr)
	require.NoError(t, err)
	assert.Equal(t, rq.RequestLine.Method, "GET")
	assert.Equal(t, rq.RequestLine.Target, "hello/")
	assert.Equal(t, rq.RequestLine.Version, "HTTP1.1")
	clength, _ := rq.Headers.Get("content-length")
	accept, _ := rq.Headers.Get("accept")
	assert.Equal(t, clength, "43")
	assert.Equal(t, accept, "text/plain, application/json")
	_, ok := rq.Headers.Get("nonexistant")
	assert.False(t, ok)
}