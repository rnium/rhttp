package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/rnium/rhttp/internal/http/headers"
)

type ParserState int8

const SEPARATOR = "\r\n"

const (
	parserInit ParserState = iota
	parserHeaders
	parserBody
	parserDone
	parserError
)

var ErrMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrMalformedFieldLine = fmt.Errorf("malformed field line")
var ErrMalformedFieldValue = fmt.Errorf("malformed field value")

type Request struct {
	RequestLine   *RequestLine
	Headers       *headers.Headers
	Body          []byte
	state         ParserState
	contentLength int
	params        Params
	query_params  Params
}

func newRequest() *Request {
	return &Request{
		state:       parserInit,
		RequestLine: &RequestLine{},
		Headers:     headers.NewHeaders(),
		Body:        nil,
	}
}

func (r *Request) done() bool {
	return r.state == parserDone || r.state == parserError
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	var err error
	sep := []byte(SEPARATOR)
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break
		}
		switch r.state {
		case parserInit:
			n, err := r.RequestLine.parseRequestLine(currentData)
			if err != nil {
				r.state = parserError
			}
			if n == 0 {
				return 0, nil // no error, just need more data
			}
			read += n
			r.state++
		case parserHeaders:
			sepIdx := bytes.Index(currentData, sep)
			if sepIdx == 0 {
				read += len(sep)
				contentLengthRaw, exists := r.Headers.Get("content-length")
				if exists {
					r.state++
					r.contentLength, err = strconv.Atoi(contentLengthRaw)
					if err != nil {
						r.state = parserError
					}
				} else {
					r.state = parserDone
				}
				continue
			} else if sepIdx == -1 {
				break outer
			}
			field_line_elems := bytes.SplitN(currentData[:sepIdx], []byte(":"), 2)
			if len(field_line_elems) != 2 {
				err = ErrMalformedFieldLine
				r.state = parserError
				continue
			}
			field_name, field_val := string(field_line_elems[0]), string(field_line_elems[1])
			err = r.Headers.Set(
				strings.TrimSpace(field_name),
				strings.TrimLeft(field_val, " "),
			)
			if err != nil {
				r.state = parserError
			}
			read += sepIdx + len(sep)
		case parserBody:
			remaining := r.contentLength - len(r.Body)
			endIdx := min(len(currentData), remaining)
			r.Body = append(r.Body, currentData[:endIdx]...)
			read += endIdx
			if r.contentLength == len(r.Body) {
				r.state = parserDone
				continue
			}
		case parserDone:
			break outer
		case parserError:
			return 0, err
		}
	}

	return read, err
}

func GetRequest(conn io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1000)
	bufIdx := 0
	for !request.done() {
		n, err := conn.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}
		n_parsed, err := request.parse(buf[:bufIdx+n])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[n_parsed:bufIdx+n])
		bufIdx = bufIdx + n - n_parsed
	}
	return request, nil
}
