package mockremotehost

import (
	"net"
	"net/url"
	"net/http"
	"strconv"
	"fmt"
)

type Server struct {
	address string
	tlsCertFile string
	tlsKeyFile string
	httpsrv http.Server
	listener net.Listener
	onRequestHandler http.HandlerFunc
}

func NewServer(host string, port int, tlsCertFile, tlsKeyFile string) *Server {
	return &Server{
		address: host + ":" + strconv.Itoa(port),
		tlsCertFile: tlsCertFile,
		tlsKeyFile: tlsKeyFile
	}
}

func (srv *Server) URL() *url.URL {
	return &url.URL{Scheme: "https", Host: srv.address}
}

func (srv *Server) BindToPort() error {
	listener, err := net.Listen("tcp", srv.address)
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
	err := srv.httpsrv.ServeTLS(srv.listener, srv.tlsCertFile, srv.tlsKeyFile)
	if err != http.ErrServerClosed {
		return fmt.Errorf("httpsrv.Serve(): %v", err)
	}

	return nil
}

func (srv *Server) OnRequest(f func(http.ResponseWriter, *http.Request)) {
	srv.onRequestHandler = http.HandlerFunc(f)
}

func (srv *Server) Stop() error{
	err := srv.httpsrv.Close()
	if err != nil { return fmt.Errorf("server.Close(): %v", err) }

	return nil
}

func (srv *Server) handlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		srv.onRequestHandler.ServeHTTP(w, req)
	})
}
