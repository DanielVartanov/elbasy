package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"io"
	"net/http/httputil"
)

type ProxyServer struct {
	URL string

	server *http.Server
}

func (proxyServer *ProxyServer) Setup() {
	proxyServer.URL = "http://localhost:8080"

	proxyServer.server = &http.Server{Addr: ":8080", Handler: proxyServer}
}

// do not run yourself (shall we have anoyter type for serving ServeHTTP interface?)
func (proxyServer *ProxyServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Received a request at Proxy")

	dump, error := httputil.DumpRequest(request, true)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println(string(dump))
	fmt.Println()

	requestCopy := *request
	requestCopy.RequestURI = ""

	fmt.Println("Making a request from Proxy to Server")

	dump, error = httputil.DumpRequestOut(&requestCopy, true)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println(string(dump))
	fmt.Println()

	responseFromServer := proxyServer.makeRequestToServer(&requestCopy)

	fmt.Println("Received a response from Server at Proxy. Relaying it to Client")

	dump, error = httputil.DumpResponse(responseFromServer, true)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Println(string(dump))
	fmt.Println()

	for headerKey, _ := range responseFromServer.Header {
		responseWriter.Header().Set(
			headerKey,
			responseFromServer.Header.Get(headerKey),
		)
	}

	responseWriter.WriteHeader(responseFromServer.StatusCode)

	_, error = io.Copy(responseWriter, responseFromServer.Body)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
}

func (proxyServer *ProxyServer) Run() {
	fmt.Println("Proxy server is running at " + proxyServer.URL)
	error := proxyServer.server.ListenAndServe()
	if error != http.ErrServerClosed {
		log.Fatal("Error in ProxyServer.Run(): ")
		log.Fatal(error)
		os.Exit(1)
	}
}

func (proxyServer *ProxyServer) Close() {
	fmt.Println("Stopping a Proxy server")
	proxyServer.server.Close()
}

func (proxyServer *ProxyServer) makeRequestToServer(request *http.Request) *http.Response {
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

	return response
}
