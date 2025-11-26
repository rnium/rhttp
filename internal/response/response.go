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
	body       []byte
	reader     io.Reader // If reader is set response will be chunked
	finished   bool
}

func NewResponse(StatusCode int, body []byte, extra_headers *headers.Headers) *Response {
	headers := getResponseHeaders(extra_headers)
	_ = headers.Set("content-length", fmt.Sprint(len(body)))
	res := &Response{
		StatusCode: StatusCode,
		Headers:    headers,
		body:       body,
	}
	return res
}

func getResponseHeaders(extra_headers *headers.Headers) *headers.Headers {
	headers := headers.GetDefaultResponseHeaders()
	if extra_headers != nil {
		extra_headers.ForEach(func(name, value string) {
			if _, exists := headers.Get(name); exists {
				_, _ = headers.Replace(name, value)
			} else {
				_ = headers.Set(name, value)
			}
		})
	}
	return headers
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
		if !res.finished {
			res.finished = true
		}
	}()

	if res.finished {
		return 0, ErrResponseClosed
	}

	// Write Status Line
	n_statusline, err := writeStatusLine(conn, res.StatusCode, request)
	if err != nil {
		return n, err
	}
	n += n_statusline

	// Write Headers
	n_fieldline, err := writeHeaders(conn, res.Headers)
	if err != nil {
		return n, err
	}
	n += n_fieldline

	// Write Body
	if res.body != nil {
		n_body, err := writeBody(conn, res.body)
		if err != nil {
			return n, err
		}
		n += n_body
	} else if res.reader != nil {
		n_body, err := res.writeChunkedBody(conn)
		if err != nil {
			return n, err
		}
		n += n_body
	}
	return
}
