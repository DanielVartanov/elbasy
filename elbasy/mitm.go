package main

import (
	"fmt"
	"net/http"
	"log"
	"net"
	"io"

	"github.com/DanielVartanov/elbasy/fakelistener"
)

type mitmConnHandler struct {
	server http.Server
	connFeeder fakelistener.ConnectionFeeder
	antiThrottler antiThrottler
}

func newMitmConnHandler() *mitmConnHandler {
	mitm := &mitmConnHandler{}

	mitm.server = http.Server{Handler: mitm}
	mitm.antiThrottler = newShopifyAntiThrottler()

	fakeListener, connFeeder := fakelistener.NewFakeListener()
	mitm.connFeeder = connFeeder

	go func() {
		err := mitm.startTLSServer(fakeListener)
		if err != nil {
			log.Fatalf("Error mitm.startTLSServer(): %v\n", err)
		}
	}()

	return mitm
}

func (mitm *mitmConnHandler) close() error {
	err := mitm.server.Close() // TODO: should be Shutdown() (or not?)
	if err != nil {
		return fmt.Errorf("ImpostorServer.server.Close(): %v", err)
	}
	return nil
}

func (mitm *mitmConnHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	mitm.antiThrottler.preventThrottling(request.Host, func(){
		remoteHostResponse, err := mitm.makeRequestToRemoteHost(request)
		if err != nil {
			log.Printf("Error on sending request to Server from Proxy: %v\n", err)
			http.Error(w, "Error sending a request to Server", 500)
			return
		}

		if remoteHostResponse.Status == "429 Too Many Requests" {
			log.Printf("Received '429 Too Many Requests' from %s\n", request.Host)
		}

		err = mitm.relayResponseToClient(w, remoteHostResponse)
		if err != nil {
			log.Printf("Error relaying a Server response from Proxy to Client: %v\n", err)
			http.Error(w, "Error relaying a Server response from Proxy to Client", 500)
		}
	})
}

// --- Private ---

func (mitm *mitmConnHandler) handleConnection(c net.Conn, _ string) error {
	mitm.connFeeder.Feed(c)
	return nil
}

func (mitm *mitmConnHandler) startTLSServer(fakeListener net.Listener) error {
	err := mitm.server.ServeTLS(fakeListener,
		"./elbasy_certificates/_wildcard.myshopify.com.pem",
		"./elbasy_certificates/_wildcard.myshopify.com-key.pem")

	if err != http.ErrServerClosed {
		return fmt.Errorf("mitm.server.ServeTLS(): %v\n", err)
	}

	return nil
}

func (mitm *mitmConnHandler) makeRequestToRemoteHost(clientReq *http.Request) (*http.Response, error) {
	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{Transport: transport}

	req := mitm.copyClientRequest(clientReq)
	response, err := client.Do(&req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do(): %v", err)
	}

	return response, nil
}

func (mitm *mitmConnHandler) copyClientRequest(clientReq *http.Request) http.Request {
	req := *clientReq
	req.RequestURI = ""
	req.URL.Scheme = "https"
	req.URL.Host = req.Host
	return req
}

func (mitm *mitmConnHandler) relayResponseToClient(w http.ResponseWriter, resp *http.Response) error {
	for headerKey, _ := range resp.Header {
		w.Header().Set(
			headerKey,
			resp.Header.Get(headerKey),
		)
	}

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("io.Copy(): %v", err)
	}

	return nil
}
