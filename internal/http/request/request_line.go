package request

import "bytes"

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
