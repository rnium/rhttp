package rhttp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"slices"
	"strings"
)

const CRLF = "\r\n"

var nonEditableResponseHeaders = []string{
	"content-length",
	"transfer-encoding",
	"connection",
	"keep-alive",
	"upgrade",
	"trailer",
	"date",
	"server",
	"via",
}

type Response struct {
	StatusCode int
	headers    *Headers
	body       []byte
	reader     io.Reader // If reader is set response will be chunked
	finished   bool
}

func NewResponse(StatusCode int, body []byte) *Response {
	headers := GetDefaultResponseHeaders()
	_ = headers.Set("content-length", fmt.Sprint(len(body)))
	res := &Response{
		StatusCode: StatusCode,
		headers:    headers,
		body:       body,
	}
	return res
}

func (r *Response) SetHeader(name, value string) error {
	if slices.Contains(nonEditableResponseHeaders, strings.ToLower(name)) {
		return ErrNonEditableHeader
	}
	var err error
	if _, exists := r.headers.Get(name); exists {
		err, _ = r.headers.Replace(name, value)
	} else {
		err = r.headers.Set(name, value)
	}
	return err
}

func writeStatusLine(conn io.Writer, StatusCode int, request *Request) (int, error) {
	statusMsg, ok := statusMessage[StatusCode]
	var statusLine string
	if ok {
		statusLine = fmt.Sprintf("%s %d %s%s", request.RequestLine.Version, StatusCode, statusMsg, CRLF)
	} else {
		statusLine = fmt.Sprintf("%s %d%s", request.RequestLine.Version, StatusCode, CRLF)
	}
	return conn.Write([]byte(statusLine))
}

func writeHeaders(conn io.Writer, headers *Headers) (int, error) {
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

func (res *Response) WriteResponse(conn io.Writer, request *Request) (n int, err error) {
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
	n_fieldline, err := writeHeaders(conn, res.headers)
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

func NewChunkedResponse(StatusCode int, reader io.Reader) *Response {
	headers := GetDefaultResponseHeaders()
	_ = headers.Set("Transfer-Encoding", "chunked")
	_ = headers.Set("Trailer", "x-content-sha256, x-content-length")
	res := &Response{
		StatusCode: StatusCode,
		headers:    headers,
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
	trailers := NewHeaders()
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
