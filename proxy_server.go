package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"net/http/httputil"
	"net"
	"time"
	"io"
)

const LEAKY_BUCKET_SIZE = 5

type ProxyServer struct {
	URL string

	server *http.Server
	throttler *throttler
}

func (proxyServer *ProxyServer) Setup() {
	proxyServer.URL = "http://localhost:8080"

	proxyServer.server = &http.Server{Addr: ":8080", Handler: proxyServer}
	proxyServer.throttler = NewThrottler(LEAKY_BUCKET_SIZE)
}

// TODO shall we have anoyter type for serving ServeHTTP interface?
func (proxyServer *ProxyServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Received a request at Proxy")

	dump, error := httputil.DumpRequest(request, false)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println(string(dump))
	fmt.Println()

	hijacker := responseWriter.(http.Hijacker)
	clientConn, _, error := hijacker.Hijack()
	if error != nil { log.Fatal(error) }

	error = clientConn.SetDeadline(time.Time{}) // Reset read/write deadlines which might have been set previously
	if error != nil { log.Fatal(error) }

	fmt.Println("Connection is hijacked. Acknowledging the proxy to client")
	_, err := clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil { log.Fatal(err) }

	serverConn, error := net.Dial("tcp", request.Host)
	if error != nil { log.Fatal(error) }
	fmt.Println("Connected to a remote Server. Starting to relay data")

	go func() {
		written, err := io.Copy(serverConn, clientConn)
		if err != nil { log.Fatal(err) }
		fmt.Println("io.Copy(serverConn, clientConn) has written", written, "bytes")
	}()

	go func() {
		written, err := io.Copy(clientConn, serverConn)
		if err != nil { log.Fatal(err) }
		fmt.Println("io.Copy(clientConn, serverConn) has written", written, "bytes")
	}()
}

func (proxyServer *ProxyServer) Run() {
	fmt.Println("Proxy server is running at " + proxyServer.URL)
	error := proxyServer.server.ListenAndServe()
	if error != http.ErrServerClosed {
		log.Fatal("Error in ProxyServer.Run(): ")
		log.Fatal(error)
		os.Exit(1)
	}
}

func (proxyServer *ProxyServer) Close() {
	fmt.Println("Stopping a Proxy server")
	proxyServer.server.Close() // TODO: should be Shutdown() (or not?)
}
