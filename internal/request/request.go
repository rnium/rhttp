package request

import (
	"bytes"
	"fmt"
	"io"
)

type ParserState int8

const SEPARATOR = "\r\n"

const (
	parserInit ParserState = iota
	parserDone
	parserError
)

var ErrMalformedRequestLine = fmt.Errorf("malformed request line")

type RequestLine struct {
	Method  string
	Target  string
	Version string
}

type Request struct {
	RequestLine *RequestLine
	state       ParserState
}

func newRequest() *Request {
	return &Request{
		state:       parserInit,
		RequestLine: &RequestLine{},
	}
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
	return len(data[:sepIdx]), nil
}

func (r *Request) done() bool {
	return r.state == parserDone || r.state == parserError
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	var err error
dance:
	for {
		switch r.state {
		case parserInit:
			n, err := r.RequestLine.parseRequestLine(data[read:])
			if err != nil {
				r.state = parserError
			}
			if n == 0 {
				return 0, nil // no error, just need more data
			}
			read += n
			r.state++
		case parserDone:
			break dance // 🕺
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
