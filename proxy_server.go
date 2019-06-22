package main

import (
	"net"
	"fmt"
	"os"
)

// proxyToServerConnection

type ProxyServer struct {
	URL string

	socketListener net.Listener
}

func (proxyServer *ProxyServer) BindToPort() {
	listener, error := net.Listen("tcp", "0.0.0.0:8080")
	proxyServer.socketListener = listener

	if error != nil {
		fmt.Println("Error on net.Listen")
		os.Exit(1)
	}

	proxyServer.URL = "http://localhost:8080"
	fmt.Println("Proxy server is running at " +
		proxyServer.socketListener.Addr().String())
}

func (proxyServer *ProxyServer) AcceptConnections() {
	for {
		clientToProxyConnection, error := proxyServer.socketListener.Accept()
		if error != nil {
			fmt.Println("Error on listener.Accept()")
			fmt.Println(error)
			os.Exit(1)
		}
		fmt.Println("Accepted a connection at Proxy from " + clientToProxyConnection.LocalAddr().String() + " | " + clientToProxyConnection.RemoteAddr().String())

	}
}

func (proxyServer *ProxyServer) Close() {
	proxyServer.socketListener.Close()
}
