package main

import (
	"fmt"
	"net/http"
	"log"
	"net/http/httputil"
	"net"
)

type FakeListener struct {
	connections chan net.Conn
	addr net.Addr
}

func NewFakeListener(addr net.Addr) *FakeListener {
	return &FakeListener{connections: make(chan net.Conn), addr: addr}
}

func (fl *FakeListener) Connect(conn net.Conn) {
	fl.connections <- conn
}

func (fl *FakeListener) Accept() (net.Conn, error) {
	conn := <-fl.connections
	return conn, nil
}

func (fl *FakeListener) Close() error {
	return nil // Should pass `shut this thing down` to another channel so that Accept() unblocks with an error
}

func (fl *FakeListener) Addr() net.Addr {
	return fl.addr
}

const LEAKY_BUCKET_SIZE = 5

type ProxyServer struct {
	listener net.Listener
	server http.Server
	fakeListener *FakeListener
}

func serveHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	dump, error := httputil.DumpRequest(request, false)
	if error != nil { log.Fatal("httputil.DumpRequest()", error)  }
	fmt.Println()
	fmt.Println(string(dump))

	responseWriter.WriteHeader(200)
	fmt.Println(responseWriter.Write([]byte("Hello, TLS world\n")))
}

func (proxyServer *ProxyServer) BindToPort() {
	proxyServer.server = http.Server{Handler: http.HandlerFunc(serveHTTP)}

	listener, err := net.Listen("tcp", ":8443")
	proxyServer.listener = listener
	if err != nil { log.Fatal("Error on net.Listen", err) }
	fmt.Println("Listening at " + listener.Addr().String())

	proxyServer.fakeListener = NewFakeListener(listener.Addr())
}

func (proxyServer *ProxyServer) AcceptConnections() {
	go func(){
		for {
			conn, err := proxyServer.listener.Accept()
			if err != nil { log.Fatal("Error on listener.Accept()", err) }
			proxyServer.fakeListener.Connect(conn)
		}
	}()

	err := proxyServer.server.ServeTLS(proxyServer.fakeListener,
		"/home/daniel/src/polite-api-proxy/localhost.crt",
		"/home/daniel/src/polite-api-proxy/localhost.key")
	if err != http.ErrServerClosed { log.Fatal("Error in http.Server.ServeTLS() ", err) }
}

func (proxyServer *ProxyServer) Run() {
	proxyServer.BindToPort()
	proxyServer.AcceptConnections()
}

func (proxyServer *ProxyServer) Close() {
	fmt.Println("Stopping a Proxy server")

	error :=  proxyServer.listener.Close()
	if error != nil { log.Fatal("Error on listener.Close()", error) }

	error = proxyServer.server.Close() // TODO: should be Shutdown() (or not?)
	if error != nil { log.Fatal("Error on server.Close()", error) }
}
