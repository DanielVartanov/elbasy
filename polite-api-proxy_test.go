package main

import (
	"testing"
	"fmt"
	"net/http"
	// "time"

	"./mock_server"
)

func TestProxyServer(t *testing.T) {
	var mockServer mock_server.MockServer
	mockServer.Start()
	defer mockServer.Close()
	fmt.Println("Mock server is running at " + mockServer.URL)

	var proxyServer ProxyServer
	proxyServer.BindToPort()
	defer proxyServer.Close()
	go proxyServer.AcceptConnections()

	clientRequestSender := ClientRequestSender{
		ServerURL: mockServer.URL,
		ProxyURL: proxyServer.URL,
	}
	clientRequestSender.SendRequest(func(_ *http.Response, responseBodyText string) {
		// Compare ALL fields of Response, including headers and the code!
		if responseBodyText != "ololo-shmololo\n" {
			t.Error("Unexpected response body:" + string(responseBodyText))
		}
	})

	clientRequestSender.WaitForAllRequests()
}
