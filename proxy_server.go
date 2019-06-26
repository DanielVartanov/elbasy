package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"io"
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

func (proxyServer *ProxyServer) readReadCloser(readCloser io.ReadCloser) string {
	defer readCloser.Close()
	fullText, error := ioutil.ReadAll(readCloser)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	return string(fullText)
}

// do not run yourself (shall we have anoyter type for that?)
func (proxyServer *ProxyServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Received a request at Proxy")

	requestCopy := *request
	requestCopy.RequestURI = ""
	// requestCopy.Body = strings.NewReader(proxyServer.readReadCloser(request.Body))

	responseFromServer := proxyServer.makeRequestToServer(&requestCopy)

	fmt.Println("Responding to the Client request from Proxy")
	fmt.Fprintf(responseWriter, responseFromServer)
}

func (proxyServer *ProxyServer) Run() {
	error := proxyServer.server.ListenAndServe()
	if error != http.ErrServerClosed {
		log.Fatal("Error in ProxyServer.Run(): ")
		log.Fatal(error)
		os.Exit(1)
	}
}

func (proxyServer *ProxyServer) Close() {
	proxyServer.server.Close()
}

func (proxyServer *ProxyServer) makeRequestToServer(request *http.Request) string {
	fmt.Println("Sending a request from Proxy to Server...")

	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
	}

	client := &http.Client{Transport: transport}

	response, error := client.Do(request)
	if error != nil {
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
