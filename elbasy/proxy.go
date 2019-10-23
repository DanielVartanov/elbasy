package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"strings"
	"strconv"
	"net/url"
)

type proxy struct {
	address string
	tlsCertFile string
	tlsKeyFile string
	server http.Server
	listener net.Listener
	mitm *mitmConnHandler
	regularProxy *regularProxyConnHandler
}

func newProxy(host string, port int, tlsCertFile, tlsKeyFile string) *proxy {
	return &proxy{
		address: host + ":" + strconv.Itoa(port),
		tlsCertFile: tlsCertFile,
		tlsKeyFile: tlsKeyFile,
		regularProxy: newRegularProxyConnHandler(),
		mitm: newMitmConnHandler(),
	}
}

func (px *proxy) url() *url.URL {
	return &url.URL{Scheme: "https", Host: px.address}
}

func (px *proxy) bindToPort() error {
	listener, err := net.Listen("tcp", px.address)
	if err != nil {
		return fmt.Errorf("net.Listen: %v", err)
	}
	px.listener = listener

	return nil
}

func (px *proxy) acceptConnections() error {
	if px.listener == nil {
		return fmt.Errorf("listener is empty, run bindToPort() first")
	}

	px.server = http.Server{Handler: px.handlerFunc()}
	err := px.server.ServeTLS(px.listener, px.tlsCertFile, px.tlsKeyFile)
	if err != http.ErrServerClosed {
		return fmt.Errorf("proxy.server.Serve(): %v", err)
	}

	return nil
}

func (px *proxy) run() error {
	err := px.bindToPort()
	if err != nil {
		return fmt.Errorf("proxy.bindToPort(): %v", err)
	}

	err = px.acceptConnections()
	if err != nil {
		return fmt.Errorf("proxy.acceptConnections(): %v", err)
	}

	return nil
}

func (px *proxy) close() error {
	err := px.server.Close()
	if err != nil {
		return fmt.Errorf("proxy.server.Close(): %v", err)
	}

	err = px.regularProxy.close()
	if err != nil {
		return fmt.Errorf("proxy.regularProxy.Close(): %v", err)
	}

	err = px.mitm.close()
	if err != nil {
		return fmt.Errorf("proxy.mitm.Close(): %v", err)
	}

	return nil
}

type connHandler interface {
	handleConnection(c net.Conn, host string) error
}

func (px proxy) handlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		clientConn, err := px.hijackConnection(responseWriter)
		if err != nil {
			log.Printf("Error hijacking connection: %v\n", err)
			http.Error(responseWriter, "Error handling request", 500)
		}

		err = px.acknowledgeProxyToClient(clientConn)
		if err != nil {
			log.Printf("Error proxy.acknowledgeProxyToClient(): %v", err)
			http.Error(responseWriter, "Error handling request", 500)
		}

		connHandler := px.chooseConnHandler(request.URL.Hostname())
		err = connHandler.handleConnection(clientConn, request.Host)
		if err != nil {
			log.Printf("Error connHandler.HandleConnection(): %v\n", err)
			http.Error(responseWriter, "Error handling request", 500)
		}
	})
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

func (px *proxy) acknowledgeProxyToClient(c net.Conn) error {
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
		return px.regularProxy
	}
}
