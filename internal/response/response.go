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
	StatusCode int
	Headers    *headers.Headers
	Body       []byte
	finished   bool
}

func NewResponse(StatusCode int, body []byte, headers *headers.Headers) *Response {
	return &Response{
		StatusCode: StatusCode,
		Headers:    headers,
		Body:       body,
	}
}

func writeStatusLine(conn io.Writer, StatusCode int, request *request.Request) (int, error) {
	statusMsg, ok := statusMessage[StatusCode]
	var statusLine string
	if ok {
		statusLine = fmt.Sprintf("%s %d %s%s", request.RequestLine.Version, StatusCode, statusMsg, CRLF)
	} else {
		statusLine = fmt.Sprintf("%s %d%s", request.RequestLine.Version, StatusCode, CRLF)
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
	if res.Headers == nil {
		res.Headers = headers.GetDefaultResponseHeaders(len(res.Body))
	} else {
		_, _ = res.Headers.Replace("content-length", fmt.Sprintf("%d", len(res.Body)))
	}
	n_statusline, err := writeStatusLine(conn, res.StatusCode, request)
	if err != nil {
		return n, err
	}
	n += n_statusline
	n_fieldline, err := writeHeaders(conn, res.Headers)
	if err != nil {
		return n, err
	}
	n += n_fieldline
	n_body, err := writeBody(conn, res.Body)
	if err != nil {
		return n, err
	}
	n += n_body
	return
}
