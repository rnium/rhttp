package response

import (
	"fmt"
	"io"

	"github.com/rnium/rhttp/internal/headers"
	"github.com/rnium/rhttp/internal/request"
)

var ErrResponseClosed = fmt.Errorf("response closed")
var ErrWriteFailed = fmt.Errorf("error while writing to connection")

const CRLF = "\r\n"

type Response struct {
	statusCode int
	headers    *headers.Headers
	body       []byte
	finished   bool
}

func NewResponse(statusCode int, body []byte, headers *headers.Headers) *Response {
	return &Response{
		statusCode: statusCode,
		headers:    headers,
		body:       body,
	}
}

func writeStatusLine(conn io.Writer, statusCode int, request *request.Request) (int, error) {
	statusMsg, ok := statusMessage[statusCode]
	var statusLine string
	if ok {
		statusLine = fmt.Sprintf("%s %d %s%s", request.RequestLine.Version, statusCode, statusMsg, CRLF)
	} else {
		statusLine = fmt.Sprintf("%s %d%s", request.RequestLine.Version, statusCode, CRLF)
	}
	return conn.Write([]byte(statusLine))
}

func writeHeaders(conn io.Writer, headers *headers.Headers) (int, error) {
	var headers_payload []byte
	headers.ForEach(func(name, value string) {
		headers_payload = fmt.Appendf(headers_payload, "%s: %s%s", name, value, CRLF)
	})
	headers_payload = fmt.Appendf(headers_payload, "%s", CRLF)
	return conn.Write(headers_payload)
}

func writeBody(conn io.Writer, p []byte) (int, error) {
	return conn.Write(p)
}

func (res *Response) WriteResponse(conn io.Writer, request *request.Request) (n int, err error) {
	defer func() {
		res.finished = true
	}()
	
	if res.finished {
		return 0, ErrResponseClosed
	}
	if res.headers == nil {
		res.headers = headers.GetDefaultResponseHeaders(len(res.body))
	} else {
		_, _ = res.headers.Replace("content-length", fmt.Sprintf("%d", len(res.body)))
	}
	n_statusline, err := writeStatusLine(conn, res.statusCode, request)
	if err != nil {
		return n, err
	}
	n += n_statusline
	n_fieldline, err := writeHeaders(conn, res.headers)
	if err != nil {
		return n, err
	}
	n += n_fieldline
	n_body, err := writeBody(conn, res.body)
	if err != nil {
		return n, err
	}
	n += n_body
	return
}
