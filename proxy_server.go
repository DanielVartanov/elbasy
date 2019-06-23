package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"io"
	"io/ioutil"
	"strings"
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
	requestBodyText := proxyServer.readReadCloser(request.Body)
	fmt.Printf("  RequestURI = %s, URL = %s, Method = %s, Body = \"%s\"\n", request.RequestURI, request.URL.String(), request.Method, requestBodyText)

	responseFromServer := proxyServer.makeRequestToServer(request.RequestURI, request.Method, requestBodyText)

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

func (proxyServer *ProxyServer) makeRequestToServer(requestURL string, requestMethod string, requestBodyText string) string {
	fmt.Println("Sending a request from Proxy to Server...")

	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
	}

	client := &http.Client{Transport: transport}

	request, error := http.NewRequest(requestMethod, requestURL, strings.NewReader(requestBodyText)) // maybe pass request.Body(io.Reader) instead of reading it in full? Use request.GetBody for debugging purposes in this case
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
