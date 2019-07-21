package artificial_listener

import (
	"net"
)

type ConnectionFeeder interface {
	Feed(net.Conn)
}

type artificialListener struct {
	connections chan net.Conn
}

func NewArtificialListener() (net.Listener, ConnectionFeeder) {
	al := artificialListener{connections: make(chan net.Conn)}
	return al, al
}

func (al artificialListener) Feed(conn net.Conn) {
	al.connections <- conn
}

func (al artificialListener) Accept() (net.Conn, error) {
	conn := <-al.connections
	return conn, nil
}

func (al artificialListener) Close() error {
	return nil // TODO: Should pass `shut this thing down` to another channel so that Accept() unblocks with an error
}

func (al artificialListener) Addr() net.Addr {
	return nil
}
