package main

import (
	"fmt"
	"net/http"
	"log"
	"net"
	"io"
)

type ElbasyServer struct {
	server http.Server
	artificialListener *ArtificialListener
	throttler *throttler
}

func NewElbasyServer() *ElbasyServer {
	es := &ElbasyServer{}

	es.server = http.Server{Handler: http.HandlerFunc(es.generateServeHTTPFunc())}
	es.artificialListener = NewArtificialListener()
	es.throttler = NewThrottler()

	go es.startTLSServer()

	return es
}

func (es *ElbasyServer) HandleConnection(conn net.Conn) {
	es.artificialListener.Connect(conn)
}

func (es *ElbasyServer) Close() {
	fmt.Println("Stopping Elbasy server")

	err := es.server.Close() // TODO: should be Shutdown() (or not?)
	if err != nil { log.Fatal("Error on server.Close()", err) }
}

// --- Private ---

func (es *ElbasyServer) startTLSServer() {
	err := es.server.ServeTLS(es.artificialListener,
		"/home/daniel/src/polite-api-proxy/elbasy_certificates/_wildcard.myshopify.com.pem",
		"/home/daniel/src/polite-api-proxy/elbasy_certificates/_wildcard.myshopify.com-key.pem")
	if err != http.ErrServerClosed { log.Fatal("Error in http.Server.ServeTLS() ", err) }
}

func (es *ElbasyServer) composeRequestToServerFromClientRequest(clientRequest *http.Request) http.Request {
	requestToServer := *clientRequest
	requestToServer.RequestURI = ""
	requestToServer.URL.Scheme = "https"
	requestToServer.URL.Host = requestToServer.Host
	return requestToServer
}

func (es *ElbasyServer) makeRequestToServer(requestFromClient *http.Request) *http.Response {
	requestToServer := es.composeRequestToServerFromClientRequest(requestFromClient)

	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{Transport: transport}

	responseFromServer, error := client.Do(&requestToServer)
	if error != nil { log.Fatal(error) }

	return responseFromServer
}

func (es *ElbasyServer) relayServerResponseToClient(responseWriter http.ResponseWriter, responseFromServer *http.Response) {
	for headerKey, _ := range responseFromServer.Header {
		responseWriter.Header().Set(
			headerKey,
			responseFromServer.Header.Get(headerKey),
		)
	}

	responseWriter.WriteHeader(responseFromServer.StatusCode)

	_, err := io.Copy(responseWriter, responseFromServer.Body)
	if err != nil { log.Fatal(err) }
}

func (es *ElbasyServer) generateServeHTTPFunc() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		es.throttler.Throttle(request.Host, func(){
			responseFromServer := es.makeRequestToServer(request)
			if responseFromServer.Status == "429 Too Many Requests" {
				log.Println("Received '429 Too Many Requests' from", request.Host)
			}
			es.relayServerResponseToClient(responseWriter, responseFromServer)
		})
	}
}
