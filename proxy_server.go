package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"io/ioutil"
)

// proxyToServerConnection

type ProxyServer struct {
	URL string

	server *http.Server
}

func (proxyServer *ProxyServer) Setup() {
	proxyServer.URL = "http://localhost:8080"

	proxyServer.server = &http.Server{Addr: ":8080", Handler: proxyServer}
}

// do not run yourself (shall we have anoyter type for that?)
func (proxyServer *ProxyServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Received a request at Proxy")
	fmt.Printf("RequestURI = %s, URL = %s\n", request.RequestURI, request.URL.String())

	responseFromServer := proxyServer.makeRequestToServer(request.RequestURI)

	fmt.Println("Responding to the Client request from Proxy")
	fmt.Fprintf(responseWriter, responseFromServer)
}

func (proxyServer *ProxyServer) Run() {
	error := proxyServer.server.ListenAndServe()
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
}

func (proxyServer *ProxyServer) Close() {
	proxyServer.server.Close()
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
