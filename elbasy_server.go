package main

import (
	"fmt"
	"net/http"
	"log"
	"net"
	"io"
	"./artificial_listener"
	"./anti_throttlers"
)

type ElbasyServer struct {
	server http.Server
	connectionFeeder artificial_listener.ConnectionFeeder
	antiThrottler anti_throttlers.AntiThrottler
}

func NewElbasyServer() *ElbasyServer {
	es := &ElbasyServer{}

	es.server = http.Server{Handler: http.HandlerFunc(es.generateServeHTTPFunc())}
	es.antiThrottler = anti_throttlers.NewShopifyAntiThrottler()

	artificialListener, connectionFeeder := artificial_listener.NewArtificialListener()
	es.connectionFeeder = connectionFeeder

	go func() {
		err := es.startTLSServer(artificialListener)
		if err != nil {
			log.Fatalf("Error ElbasyServer.startTLSServer(): %v", err)
		}
	}()

	return es
}

func (es *ElbasyServer) HandleConnection(conn net.Conn) {
	es.connectionFeeder.Feed(conn)
}

func (es *ElbasyServer) Close() error {
	err := es.server.Close() // TODO: should be Shutdown() (or not?)
	if err != nil {
		return fmt.Errorf("ElbasyServer.server.Close(): %v", err)
	} else {
		return nil
	}
}

// --- Private ---

func (es *ElbasyServer) startTLSServer(artificialListener net.Listener) error {
	err := es.server.ServeTLS(artificialListener,
		"/home/daniel/src/polite-api-proxy/elbasy_certificates/_wildcard.myshopify.com.pem",
		"/home/daniel/src/polite-api-proxy/elbasy_certificates/_wildcard.myshopify.com-key.pem")

	if err != http.ErrServerClosed {
		return fmt.Errorf("ElbasyServer.server.ServeTLS(): %v", err)
	} else {
		return nil
	}
}

func (es *ElbasyServer) composeRequestToServerFromClientRequest(clientRequest *http.Request) http.Request {
	requestToServer := *clientRequest
	requestToServer.RequestURI = ""
	requestToServer.URL.Scheme = "https"
	requestToServer.URL.Host = requestToServer.Host
	return requestToServer
}

func (es *ElbasyServer) makeRequestToServer(requestFromClient *http.Request) (*http.Response, error) {
	requestToServer := es.composeRequestToServerFromClientRequest(requestFromClient)

	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{Transport: transport}

	responseFromServer, err := client.Do(&requestToServer)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do(): %v", err)
	} else {
		return responseFromServer, nil
	}
}

func (es *ElbasyServer) relayServerResponseToClient(responseWriter http.ResponseWriter, responseFromServer *http.Response) error {
	for headerKey, _ := range responseFromServer.Header {
		responseWriter.Header().Set(
			headerKey,
			responseFromServer.Header.Get(headerKey),
		)
	}

	responseWriter.WriteHeader(responseFromServer.StatusCode)

	_, err := io.Copy(responseWriter, responseFromServer.Body)
	if err != nil {
		return fmt.Errorf("io.Copy(): %v", err)
	} else {
		return nil
	}
}

func (es *ElbasyServer) generateServeHTTPFunc() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		es.antiThrottler.PreventThrottling(request.Host, func(){
			responseFromServer, err := es.makeRequestToServer(request)
			if err != nil {
				log.Printf("Error on sending request to Server from Proxy: %v", err)
				http.Error(responseWriter, "Error sending a request to Server", 500)
				return
			}

			if responseFromServer.Status == "429 Too Many Requests" {
				log.Print("Received '429 Too Many Requests' from ", request.Host)
			}

			err = es.relayServerResponseToClient(responseWriter, responseFromServer)
			if err != nil {
				log.Printf("Error relaying a Server response from Proxy to Client: %v", err)
				http.Error(responseWriter, "Error relaying a Server response from Proxy to Client", 500)
			}
		})
	}
}
