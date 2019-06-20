package main

// How to exit more gracefully instead of os.Exit?

import (
	"fmt"
	"net"
	"os"
	// "time"
)

func main() {
	listener, error := net.Listen("tcp", "0.0.0.0:8080")
	if error != nil {
		fmt.Println("Error on net.Listen")
		os.Exit(1)
	}
	fmt.Println("Listening the port 8080")
	fmt.Println("listener.Addr() == " + listener.Addr().String())

	for {
		client_connection, error := listener.Accept()
		if error != nil {
			fmt.Println("Error on listener.Accept()")
			os.Exit(1)
		}
		fmt.Println("Accepted a connection from " + client_connection.LocalAddr().String() + " | " + client_connection.RemoteAddr().String())
	}
}
