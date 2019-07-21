package impostor_server

import (
	"fmt"
	"net/http"
	"log"
	"net"
	"io"
	"./artificial_listener"
	"./anti_throttlers"
)

type ImpostorServer struct {
	server http.Server
	connectionFeeder artificial_listener.ConnectionFeeder
	antiThrottler anti_throttlers.AntiThrottler
}

func NewImpostorServer() *ImpostorServer {
	is := &ImpostorServer{}

	is.server = http.Server{Handler: http.HandlerFunc(is.generateServeHTTPFunc())}
	is.antiThrottler = anti_throttlers.NewShopifyAntiThrottler()

	artificialListener, connectionFeeder := artificial_listener.NewArtificialListener()
	is.connectionFeeder = connectionFeeder

	go func() {
		err := is.startTLSServer(artificialListener)
		if err != nil {
			log.Fatalf("Error ImpostorServer.startTLSServer(): %v", err)
		}
	}()

	return is
}

func (is *ImpostorServer) HandleConnection(conn net.Conn) {
	is.connectionFeeder.Feed(conn)
}

func (is *ImpostorServer) Close() error {
	err := is.server.Close() // TODO: should be Shutdown() (or not?)
	if err != nil {
		return fmt.Errorf("ImpostorServer.server.Close(): %v", err)
	} else {
		return nil
	}
}

// --- Private ---

func (is *ImpostorServer) startTLSServer(artificialListener net.Listener) error {
	err := is.server.ServeTLS(artificialListener,
		"./elbasy_certificates/_wildcard.myshopify.com.pem",
		"./elbasy_certificates/_wildcard.myshopify.com-key.pem")

	if err != http.ErrServerClosed {
		return fmt.Errorf("ImpostorServer.server.ServeTLS(): %v", err)
	} else {
		return nil
	}
}

func (is *ImpostorServer) composeRequestToServerFromClientRequest(clientRequest *http.Request) http.Request {
	requestToServer := *clientRequest
	requestToServer.RequestURI = ""
	requestToServer.URL.Scheme = "https"
	requestToServer.URL.Host = requestToServer.Host
	return requestToServer
}

func (is *ImpostorServer) makeRequestToServer(requestFromClient *http.Request) (*http.Response, error) {
	requestToServer := is.composeRequestToServerFromClientRequest(requestFromClient)

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

func (is *ImpostorServer) relayServerResponseToClient(responseWriter http.ResponseWriter, responseFromServer *http.Response) error {
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

func (is *ImpostorServer) generateServeHTTPFunc() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		is.antiThrottler.PreventThrottling(request.Host, func(){
			responseFromServer, err := is.makeRequestToServer(request)
			if err != nil {
				log.Printf("Error on sending request to Server from Proxy: %v", err)
				http.Error(responseWriter, "Error sending a request to Server", 500)
				return
			}

			if responseFromServer.Status == "429 Too Many Requests" {
				log.Print("Received '429 Too Many Requests' from ", request.Host)
			}

			err = is.relayServerResponseToClient(responseWriter, responseFromServer)
			if err != nil {
				log.Printf("Error relaying a Server response from Proxy to Client: %v", err)
				http.Error(responseWriter, "Error relaying a Server response from Proxy to Client", 500)
			}
		})
	}
}
