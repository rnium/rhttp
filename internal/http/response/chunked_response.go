package response

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/rnium/rhttp/internal/http/headers"
)

func NewChunkedResponse(StatusCode int, reader io.Reader, extra_headers *headers.Headers) *Response {
	headers := getResponseHeaders(extra_headers)
	_ = headers.Set("Transfer-Encoding", "chunked")
	_ = headers.Set("Trailer", "x-content-sha256, x-content-length")
	res := &Response{
		StatusCode: StatusCode,
		Headers:    headers,
		reader:     reader,
	}
	return res
}

func (res *Response) writeChunkedBody(conn io.Writer) (n int, err error) {
	var body []byte
	var errored bool
	for {
		buf := make([]byte, 512)
		n_read, err := res.reader.Read(buf)
		if err != nil {
			break
		}
		if n_read > 0 {
			chunk := buf[:n_read]
			body = append(body, chunk...)
			n_body, err := writeBody(conn, fmt.Appendf(nil, "%x%s", n_read, CRLF))
			if err != nil {
				errored = true
				break
			}
			n += n_body
			n_body, err = writeBody(conn, fmt.Append(chunk, CRLF))
			if err != nil {
				errored = true
				break
			}
			n += n_body
		}
	}
	if errored {
		return
	}
	n_body, err := writeBody(conn, fmt.Appendf(nil, "0%s", CRLF))
	if err != nil {
		return
	}
	n += n_body
	trailers := headers.NewHeaders()
	digest := sha256.Sum256(body)
	_ = trailers.Set("x-content-sha256", hex.EncodeToString(digest[:]))
	_ = trailers.Set("x-content-length", fmt.Sprint(len(body)))
	n_trailer, err := writeHeaders(conn, trailers)
	if err != nil {
		return n, err
	}
	n += n_trailer
	return
}
