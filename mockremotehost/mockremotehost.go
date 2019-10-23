package mockremotehost

import (
	"net"
	"net/url"
	"net/http"
	"strconv"
	"fmt"
)

type Server struct {
	Port int

	httpsrv http.Server
	listener net.Listener
	onRequestHandler http.HandlerFunc
}

func (srv *Server) URL() *url.URL {
	return &url.URL{Scheme: "http", Host: "localhost:" + strconv.Itoa(srv.Port)}
}

func (srv *Server) BindToPort() error {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(srv.Port))
	if err != nil {
		return fmt.Errorf("net.Listen: %v", err)
	}
	srv.listener = listener
	return nil
}

func (srv *Server) AcceptConnections() error {
	if srv.listener == nil {
		return fmt.Errorf("listener is empty, run BindToPort() first")
	}

	srv.httpsrv = http.Server{Handler: srv.handlerFunc()}
	err := srv.httpsrv.Serve(srv.listener)
	if err != http.ErrServerClosed {
		return fmt.Errorf("httpsrv.Serve(): %v", err)
	}

	return nil
}

func (srv *Server) OnRequest(f func(http.ResponseWriter, *http.Request)) {
	srv.onRequestHandler = http.HandlerFunc(f)
}

func (srv *Server) handlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		srv.onRequestHandler.ServeHTTP(w, req)
	})
}
