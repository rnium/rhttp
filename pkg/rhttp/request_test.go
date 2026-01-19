package rhttp

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
		bytesPerRead: 3,
		pos: 0,
	}
}



func TestRequest(t *testing.T) {	
	cr := newChunkedReader([]byte("GET hello/ HTTP1.1\r\n\r\n"))
	rq, err := GetRequest(cr)
	require.NoError(t, err)
	assert.Equal(t, rq.RequestLine.Method, "GET")
	assert.Equal(t, rq.RequestLine.Target, "hello/")
	assert.Equal(t, rq.RequestLine.Version, "HTTP1.1")

}


func TestRequestHeaders(t *testing.T) {
	cr := newChunkedReader([]byte("GET hello/ HTTP1.1\r\ncontent-length: 5\r\naccept: text/plain\r\naccept: application/json\r\n\r\nHello\r\n"))
	rq, err := GetRequest(cr)
	require.NoError(t, err)
	
	clength, _ := rq.Headers.Get("content-length")
	accept, _ := rq.Headers.Get("accept")
	assert.Equal(t, clength, "5")
	assert.Equal(t, accept, "text/plain, application/json")
	_, ok := rq.Headers.Get("nonexistant")
	assert.False(t, ok)
	assert.Equal(t, string(rq.Body), "Hello")
	cr = newChunkedReader([]byte("GET hello/ HTTP1.1\r\nhost: localhost:80\r\n\r\n"))
	_, err = GetRequest(cr)
	require.NoError(t, err)
	cr = newChunkedReader([]byte("GET hello/ HTTP1.1\r\nhost localhost:80\r\n\r\n"))
	_, err = GetRequest(cr)
	assert.ErrorIs(t, err, ErrInvalidToken)

	cr = newChunkedReader([]byte("GET hello/ HTTP1.1\r\nhost\r\n\r\n"))
	_, err = GetRequest(cr)
	assert.ErrorIs(t, err, ErrMalformedFieldLine)
}
