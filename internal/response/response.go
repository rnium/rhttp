package response

import (
	"fmt"
	"io"

	"github.com/rnium/rhttp/internal/headers"
	"github.com/rnium/rhttp/internal/request"
)

const CRLF = "\r\n"

type Response struct {
	request *request.Request
}

func NewResponse(req *request.Request) *Response {
	return &Response{
		request: req,
	}
}


func (res *Response) WriteStatusLine(conn io.Writer, statusCode int) (int, error) {
	statusMsg, ok := statusMessage[statusCode]
	var statusLine string
	if ok {
		statusLine = fmt.Sprintf("%s %d %s%s", res.request.RequestLine.Version, statusCode, statusMsg, CRLF)
	} else {
		statusLine = fmt.Sprintf("%s %d%s", res.request.RequestLine.Version, statusCode, CRLF)
	}
	return conn.Write([]byte(statusLine))
}


func (res *Response) WriteResponseHeaders(conn io.Writer, headers *headers.Headers) (int, error) {
	var headers_payload []byte
	headers.ForEach(func(name, value string) {
		headers_payload = fmt.Appendf(headers_payload, "%s: %s%s", name, value, CRLF)
	})
	headers_payload = fmt.Appendf(headers_payload, "%s", CRLF)
	return conn.Write(headers_payload)
}