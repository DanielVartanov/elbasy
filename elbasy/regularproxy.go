package main

import (
	"fmt"
	"net"
	"io"
	"log"
)

type regularProxyConnHandler struct {}

func newRegularProxyConnHandler() *regularProxyConnHandler {
	return &regularProxyConnHandler{}
}

func (_ *regularProxyConnHandler) handleConnection(clientConn net.Conn, remoteHost string) error {
	serverConn, err := net.Dial("tcp", remoteHost)
	if err != nil {
		return fmt.Errorf("net.Dial: %v", err)
	}

	go func() {
		_, err = io.Copy(serverConn, clientConn)
		if err != nil {
			log.Printf("Error regularConnHandler.handleConnection(): io.Copy(serverConn, clientConn): %v\n", err)
		}

	}()

	go func() {
		_, err = io.Copy(clientConn, serverConn)
		if err != nil {
			log.Printf("Error regularConnHandler.handleConnection(): io.Copy(clientConn, serverConn): %v\n", err)
		}
	}()

	return nil
}

func (_ *regularProxyConnHandler) close() error {
	return nil
}
