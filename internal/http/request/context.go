package request

import (
	"net"

	"github.com/rnium/rhttp/internal/http/form"
)

func (req *Request) RemoteAddr() net.Addr {
	return (*req.conn).RemoteAddr()
}

func (req *Request) LocalAddr() net.Addr {
	return (*req.conn).LocalAddr()
}

func (req *Request) Conn() net.Conn {
	return *req.conn
}

func (req *Request) FormData() (*form.FormData, error) {
	if req == nil {
		return nil, form.ErrNoFormData
	}
	ctype, _ := req.Headers.Get("content-type")
	return form.GetFormData(ctype, req.Body)
}
