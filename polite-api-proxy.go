package main

// How to exit more gracefully instead of os.Exit?

import (
	"fmt"
	"net"
	"os"
	// "time"
)

func RunProxyServer() net.Listener {
	listener, error := net.Listen("tcp", "0.0.0.0:8080")
	if error != nil {
		fmt.Println("Error on net.Listen")
		os.Exit(1)
	}
	fmt.Println("Proxy server is running at " + listener.Addr().String())

	go func() {
		for {
			clientConnection, error := listener.Accept()
			if error != nil {
				fmt.Println("Error on listener.Accept()")
				fmt.Println(error)
				os.Exit(1)
			}
			fmt.Println("Accepted a connection at Proxy from " + clientConnection.LocalAddr().String() + " | " + clientConnection.RemoteAddr().String())

		}
	}()
	return listener
}

func main() {
	RunProxyServer()
}
