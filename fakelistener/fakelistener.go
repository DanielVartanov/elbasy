package fakelistener

import (
	"net"
)

type ConnectionFeeder interface {
	Feed(net.Conn)
}

type fakeListener struct {
	connections chan net.Conn
}

func NewFakeListener() (net.Listener, ConnectionFeeder) {
	fl := fakeListener{connections: make(chan net.Conn)}
	return fl, fl
}

func (fl fakeListener) Feed(conn net.Conn) {
	fl.connections <- conn
}

func (fl fakeListener) Accept() (net.Conn, error) {
	conn := <-fl.connections
	return conn, nil
}

func (fl fakeListener) Close() error {
	return nil // TODO: Should pass `shut this thing down` to another channel so that Accept() unblocks with an error
}

func (fl fakeListener) Addr() net.Addr {
	return nil
}
