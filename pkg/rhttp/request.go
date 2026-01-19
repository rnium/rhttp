package rhttp

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
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

type RequestLine struct {
	Method  string
	Target  string
	Version string
}

func (rl *RequestLine) parseRequestLine(data []byte) (int, error) {
	sep := []byte(SEPARATOR)
	sepIdx := bytes.Index(data, sep)
	if sepIdx == -1 {
		return 0, nil
	}
	elements_data := bytes.Split(data[:sepIdx], []byte(" "))
	if len(elements_data) != 3 {
		return 0, ErrMalformedRequestLine
	}
	rl.Method = string(elements_data[0])
	rl.Target = string(elements_data[1])
	rl.Version = string(elements_data[2])
	return sepIdx + len(sep), nil
}

type Params map[string]string

func NewParams() Params {
	return make(Params)
}

type Request struct {
	RequestLine   *RequestLine
	Headers       *Headers
	Body          []byte
	state         ParserState
	contentLength int
	conn          *net.Conn
	params        Params
	query_params  Params
}

func newRequest(conn *net.Conn) *Request {
	return &Request{
		state:       parserInit,
		RequestLine: &RequestLine{},
		Headers:     NewHeaders(),
		Body:        nil,
		conn:        conn,
	}
}

func (r *Request) SetAllParams(params Params, query_params Params) {
	r.params = params
	r.query_params = query_params
}

func (r *Request) Param(name string) (string, bool) {
	value, ok := r.params[name]
	return value, ok
}

func (r *Request) QParam(name string) (string, bool) {
	value, ok := r.query_params[name]
	return value, ok
}

func (r *Request) QParamForEach(f func(name, value string)) {
	for name, val := range r.query_params {
		f(name, val)
	}
}

func (req *Request) RemoteAddr() net.Addr {
	return (*req.conn).RemoteAddr()
}

func (req *Request) LocalAddr() net.Addr {
	return (*req.conn).LocalAddr()
}

func (req *Request) Conn() net.Conn {
	return *req.conn
}

func (req *Request) FormData() (*FormData, error) {
	if req == nil {
		return nil, ErrNoFormData
	}
	ctype, _ := req.Headers.Get("content-type")
	return GetFormData(ctype, req.Body)
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
					r.contentLength, err = strconv.Atoi(contentLengthRaw)
					if err != nil {
						r.state = parserError
					}
					if r.contentLength > 0 {
						r.state = parserBody
					}
				}
				if !exists || r.contentLength == 0 {
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
	netConn, _ := conn.(net.Conn)
	request := newRequest(&netConn)
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
