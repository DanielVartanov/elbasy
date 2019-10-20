package fakelistener

import (
	"testing"
	"net"
	"net/http"
)

func TestItFeedsAConnection(t *testing.T) {
	_, wantedConn := net.Pipe()

	listener, feeder := NewFakeListener()

	go func(){ feeder.Feed(wantedConn) }()
	gotConn, err := listener.Accept()
	if err != nil {
		t.Errorf("Error in listener.Accept(): %v", err)
	}

	if gotConn != wantedConn {
		t.Error("gotConn != wantedConn")
	}
}

func TestAsItIsIntendedToBeUsedInReality(t *testing.T) {
	clientConn, serverConn := net.Pipe()

	listener, feeder := NewFakeListener()

	wantMsg := "3.141592"
	waitCh := make(chan int)
	go func(){
		server := http.Server{Handler: http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request){
			_, receivedConn, err := responseWriter.(http.Hijacker).Hijack()
			if err != nil {
				t.Errorf("Error in hijacker.Hijack(): %v", err)
			}

			buf := make([]byte, 8)
			_, err = receivedConn.Read(buf)
			if err != nil {
				t.Errorf("Error in receivedConn.Read(): %v", err)
			}
			gotMsg := string(buf)

			if gotMsg != wantMsg {
				t.Errorf("gotMsg != wantMsg: \"%s\" != \"%s\"", gotMsg, wantMsg)
			}

			waitCh <- 1
		})}

		err := server.Serve(listener)
		if err != http.ErrServerClosed {
			t.Errorf("Error in server.Serve(): %v", err)
		}
	}()

	feeder.Feed(serverConn)
	_, err := clientConn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n" + wantMsg))
	if err != nil {
		t.Errorf("Error in clientConn.Write(): %v", err)
	}

	<-waitCh
}
