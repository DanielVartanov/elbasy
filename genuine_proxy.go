package main

import (
	"net"
	"log"
	"io"
)

type GenuineProxy struct {}

func (gp *GenuineProxy) HandleConnection(clientConn net.Conn, remoteHost string) {
	serverConn, err := net.Dial("tcp", remoteHost)
	if err != nil { log.Fatal(err) }

	go io.Copy(serverConn, clientConn)
	go io.Copy(clientConn, serverConn)
}
