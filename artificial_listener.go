package main

import (
	"net"
)

type ArtificialListener struct {
	connections chan net.Conn
}

func NewArtificialListener() *ArtificialListener {
	return &ArtificialListener{connections: make(chan net.Conn)}
}

func (al *ArtificialListener) Connect(conn net.Conn) {
	al.connections <- conn
}

func (al *ArtificialListener) Accept() (net.Conn, error) {
	conn := <-al.connections
	return conn, nil
}

func (al *ArtificialListener) Close() error {
	return nil // Should pass `shut this thing down` to another channel so that Accept() unblocks with an error
}

func (al *ArtificialListener) Addr() net.Addr {
	return nil
}
