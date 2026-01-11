package request

import "net"

func (req *Request) RemoteAddr() net.Addr {
	return (*req.conn).RemoteAddr()
}

func (req *Request) LocalAddr() net.Addr {
	return (*req.conn).LocalAddr()
}

func (req *Request) Conn() net.Conn {
	return *req.conn
}