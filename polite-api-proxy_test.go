package main

import (
	"testing"
	"os"
	// "net"
	"net/http"
	"fmt"
	"io/ioutil"
	"time"
	"net/url"
	"./mock_server"
)

/*

proxyServer = proxy_server.NewProxyServer(...)
defer proxyServer.Close()

*/

func TestProxyServer(t *testing.T) {
	var mockServer mock_server.MockServer
	mockServer.Start()
	defer mockServer.Close()
	fmt.Println("Mock server is running at " + mockServer.URL)

	var proxyServer ProxyServer
	proxyServer.BindToPort()
	defer proxyServer.Close()
	go proxyServer.AcceptConnections()

	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8080") // politeAPIProxy.ProxyURL()
	proxyUrl, error := url.Parse("http://localhost:8080") // Why env var not working?
	if error != nil {
		t.Error("Error parsing proxy URL")
	}
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	fmt.Println("Initiating a client request...")
	response, error := http.Get(mockServer.URL)
	if error != nil {
		t.Error(error)
	}

	responseBodyText, error := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if error != nil {
		t.Error(error)
	}
	// Compare ALL fields of Response, including headers and the code!
	if string(responseBodyText) != "ololo-shmololo" {
		t.Error("Unexpected response body:" + string(responseBodyText))
	}

	time.Sleep(1 * time.Second)
}
