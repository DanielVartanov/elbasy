package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"strings"
)

type ProxyServer struct {
	elbasyServer *ElbasyServer
	server http.Server
	listener net.Listener
	genuineProxy GenuineProxy
}

func (ps *ProxyServer) BindToPort() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("net.Listen: %v", err)
	}
	ps.listener = listener
	log.Print("Listening at " + listener.Addr().String())

	ps.server = http.Server{Handler: http.HandlerFunc(ps.generateServeHTTPFunc())}
	ps.elbasyServer = NewElbasyServer()

	return nil
}

func (ps *ProxyServer) AcceptConnections() error {
	err := ps.server.Serve(ps.listener)
	if err != http.ErrServerClosed {
		return fmt.Errorf("ProxyServer.server.Serve(): %v", err)
	}
	return nil
}

func (ps *ProxyServer) Run() error {
	err := ps.BindToPort()
	if err != nil {
		return fmt.Errorf("ProxyServer.BindToPort(): %v", err)
	}

	err = ps.AcceptConnections()
	if err != nil {
		return fmt.Errorf("ProxyServer.AcceptConnections(): %v", err)
	}

	return nil
}

func (ps *ProxyServer) Close() error {
	fmt.Println("Stopping a Proxy server")

	err :=  ps.listener.Close()
	if err != nil {
		return fmt.Errorf("ProxyServer.listener.Close(): %v", err)
	}

	err = ps.server.Close()
	if err != nil {
		return fmt.Errorf("ProxyServer.server.Close(): %v", err)
	}

	err = ps.elbasyServer.Close()
	if err != nil {
		return fmt.Errorf("ProxyServer.elbasyServer.Close(): %v", err)
	}

	return nil
}

// --- Private ---

func (ps *ProxyServer) hijackConnection(responseWriter http.ResponseWriter) (net.Conn, error) {
	hijacker := responseWriter.(http.Hijacker)
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("hijacker.Hijack(): %v", err)
	}

	err = clientConn.SetDeadline(time.Time{}) // Reset read/write deadlines which might have been set previously
	if err != nil {
		return nil, fmt.Errorf("clientConn.SerDeadline(): %v", err)
	}

	return clientConn, nil
}

func (ps *ProxyServer) acknowledgeProxyToClient(clientConn net.Conn) error {
	_, err := clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		return fmt.Errorf("clientConn.Write(): %v", err)
	}
	return nil
}

func (ps *ProxyServer) isHostShopify(host string) bool {
	return strings.HasSuffix(host, ".myshopify.com")
}

func (ps *ProxyServer) generateServeHTTPFunc() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		clientConn, err := ps.hijackConnection(responseWriter)
		if err != nil {
			log.Printf("Error hijacking connection: %v", err)
			http.Error(responseWriter, "Error handling client connection", 500)
		}

		err = ps.acknowledgeProxyToClient(clientConn)
		if err != nil {
			log.Printf("Error ProxyServer.acknowledgeProxyToClient(): %v", err)
			http.Error(responseWriter, "Error handling client connection", 500)
		}

		if ps.isHostShopify(request.URL.Hostname()) {
			ps.elbasyServer.HandleConnection(clientConn)
		} else {
			err = ps.genuineProxy.HandleConnection(clientConn, request.Host)
			if err != nil {
				log.Printf("Error ProxyServer.genuineProxy.HandleConnection(): %v", err)
				http.Error(responseWriter, "Error handling client connection", 500)
			}
		}
	}
}
