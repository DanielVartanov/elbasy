package main

import (
	"net"
	"fmt"
	"os"
	"net/http"
	"log"
	"io/ioutil"
)

// proxyToServerConnection

type ProxyServer struct {
	URL string

	socketListener net.Listener
}

func (proxyServer *ProxyServer) BindToPort() {
	proxyServer.URL = "http://localhost:8080"

	/*
	listener, error := net.Listen("tcp", "0.0.0.0:8080")
	if error != nil {
		fmt.Println("Error on net.Listen")
		os.Exit(1)
	}
	proxyServer.socketListener = listener
	proxyServer.URL = "http://localhost:8080"

	fmt.Println("Proxy server is running at " +
		proxyServer.socketListener.Addr().String())

*/
}

func (proxyServer *ProxyServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Received a request at Proxy")
	fmt.Printf("RequestURI = %s, URL = %s\n", request.RequestURI, request.URL.String())

	responseFromServer := proxyServer.makeRequestToServer(request.RequestURI)

	fmt.Println("Responding to the Client request from Proxy")
	fmt.Fprintf(responseWriter, responseFromServer)
}

func (proxyServer *ProxyServer) AcceptConnections() {
	error := http.ListenAndServe(":8080", proxyServer)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}

	/*
	for {
		clientToProxyConnection, error := proxyServer.socketListener.Accept()
		if error != nil {
			fmt.Println("Error on listener.Accept()")
			fmt.Println(error)
			os.Exit(1)
		}
		fmt.Println("Accepted a connection at Proxy from " + clientToProxyConnection.LocalAddr().String() + " | " + clientToProxyConnection.RemoteAddr().String())

		go proxyServer.handleConnection(&clientToProxyConnection)

	}
*/
}

/*
func (proxyServer *ProxyServer) handleConnection(clientToProxyConnection *Conn) {

}
*/

func (proxyServer *ProxyServer) Close() {
	// proxyServer.socketListener.Close()
}

func (proxyServer *ProxyServer) makeRequestToServer(requestURL string) string {
	fmt.Println("Sending a request from Proxy to Server...")

	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
	}

	client := &http.Client{Transport: transport}

	request, error := http.NewRequest("GET", requestURL, nil)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(error)
		os.Exit(1)
	}

	responseBody, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
	defer response.Body.Close()

	return string(responseBody)
}
