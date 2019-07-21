package genuine_proxy

import (
	"net"
	"log"
	"io"
	"fmt"
)

type GenuineProxy struct {}

func (gp *GenuineProxy) HandleConnection(clientConn net.Conn, remoteHost string) error {
	serverConn, err := net.Dial("tcp", remoteHost)
	if err != nil {
		return fmt.Errorf("net.Dial: %v", err)
	}

	go func() {
		_, err = io.Copy(serverConn, clientConn)
		if err != nil {
			log.Printf("Error GenuineProxy.HandleConnection(): io.Copy(serverConn, clientConn): %v", err)
		}

	}()

	go func() {
		_, err = io.Copy(clientConn, serverConn)
		if err != nil {
			log.Printf("Error GenuineProxy.HandleConnection(): io.Copy(clientConn, serverConn): %v", err)
		}
	}()

	return nil
}
