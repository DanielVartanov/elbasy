package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"strings"
)

type proxy struct {
	server http.Server
	listener net.Listener

	mitm *mitmConnHandler
	regular *regularConnHandler
}

func newProxy() *proxy {
	var px proxy
	px.regular = newRegularConnHandler()
	px.mitm = newMitmConnHandler()
	px.server = http.Server{Handler: px}
	return &px
}

func (px *proxy) bindToPort() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("net.Listen: %v", err)
	}
	px.listener = listener
	log.Print("Listening at " + px.listener.Addr().String())
	return nil
}

func (px *proxy) acceptConnections() error {
	err := px.server.Serve(px.listener)
	if err != http.ErrServerClosed {
		return fmt.Errorf("proxy.server.Serve(): %v", err)
	}
	return nil
}

func (px *proxy) run() error {
	err := px.bindToPort()
	if err != nil {
		return fmt.Errorf("proxy.BindToPort(): %v", err)
	}

	err = px.acceptConnections()
	if err != nil {
		return fmt.Errorf("proxy.AcceptConnections(): %v", err)
	}

	return nil
}

func (px *proxy) close() error {
	log.Println("Stopping a proxy server")

	err :=  px.listener.Close()
	if err != nil {
		return fmt.Errorf("proxy.listener.Close(): %v", err)
	}

	err = px.server.Close()
	if err != nil {
		return fmt.Errorf("proxy.server.Close(): %v", err)
	}

	err = px.regular.close()
	if err != nil {
		return fmt.Errorf("proxy.regular.Close(): %v", err)
	}

	err = px.mitm.close()
	if err != nil {
		return fmt.Errorf("proxy.mitm.Close(): %v", err)
	}

	return nil
}

// --- Private ---

type connHandler interface {
	handleConnection(c net.Conn, host string) error
}

func (px proxy) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	clientConn, err := px.hijackConnection(responseWriter)
	if err != nil {
		log.Printf("Error hijacking connection: %v\n", err)
		http.Error(responseWriter, "Error handling request", 500)
	}

	err = px.acknowledgeproxyToClient(clientConn)
	if err != nil {
		log.Printf("Error proxy.acknowledgeproxyToClient(): %v", err)
		http.Error(responseWriter, "Error handling request", 500)
	}

	connHandler := px.chooseConnHandler(request.URL.Hostname())
	err = connHandler.handleConnection(clientConn, request.Host)
	if err != nil {
		log.Printf("Error connHandler.HandleConnection(): %v\n", err)
		http.Error(responseWriter, "Error handling request", 500)
	}
}

func (px *proxy) hijackConnection(responseWriter http.ResponseWriter) (net.Conn, error) {
	hijacker := responseWriter.(http.Hijacker)
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("hijacker.Hijack(): %v", err)
	}

	err = clientConn.SetDeadline(time.Time{}) // Reset read/write deadlines which might have been set previously
	if err != nil {
		return nil, fmt.Errorf("clientConn.SetDeadline(): %v", err)
	}

	return clientConn, nil
}

func (px *proxy) acknowledgeproxyToClient(c net.Conn) error {
	_, err := c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		return fmt.Errorf("net.Conn.Write(): %v", err)
	}
	return nil
}

func (px *proxy) chooseConnHandler(host string) connHandler {
	if strings.HasSuffix(host, ".myshopify.com") {
		return px.mitm
	} else {
		return px.regular
	}
}
