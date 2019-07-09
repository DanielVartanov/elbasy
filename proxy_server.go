package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type ProxyServer struct {
	elbasyServer *ElbasyServer
	server http.Server
	listener net.Listener
}

func (ps *ProxyServer) BindToPort() {
	listener, err := net.Listen("tcp", ":8080")
	ps.listener = listener
	if err != nil { log.Fatal("Error on net.Listen", err) }

	fmt.Println("Listening at " + listener.Addr().String())

	ps.server = http.Server{Handler: http.HandlerFunc(ps.generateServeHTTPFunc())}
	ps.elbasyServer = NewElbasyServer()
}

func (ps *ProxyServer) AcceptConnections() {
	ps.server.Serve(ps.listener)
}

func (ps *ProxyServer) Run() {
	ps.BindToPort()
	ps.AcceptConnections()
}

func (ps *ProxyServer) Close() {
	fmt.Println("Stopping a Proxy server")

	error :=  ps.listener.Close()
	if error != nil { log.Fatal("Error on listener.Close()", error) }

	error = ps.server.Close() // TODO: should be Shutdown() (or not?)
	if error != nil { log.Fatal("Error on server.Close()", error) }

	ps.elbasyServer.Close()
}

// --- Private ---

func (ps *ProxyServer) hijackConnection(responseWriter http.ResponseWriter) net.Conn {
	hijacker := responseWriter.(http.Hijacker)
	clientConn, _, error := hijacker.Hijack()
	if error != nil { log.Fatal(error) }

	error = clientConn.SetDeadline(time.Time{}) // Reset read/write deadlines which might have been set previously
	if error != nil { log.Fatal(error) }

	return clientConn
}

func (ps *ProxyServer) acknowledgeProxyToClient(clientConn net.Conn) {
	_, err := clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil { log.Fatal(err) }
}

func (ps *ProxyServer) generateServeHTTPFunc() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		clientConn := ps.hijackConnection(responseWriter)

		ps.acknowledgeProxyToClient(clientConn)

		ps.elbasyServer.HandleConnection(clientConn)
	}
}
