package main

import (
	"testing"
	"fmt"
	"net/http"
	// "time"

	"./mock_server"
)

func startMockServer() *mock_server.MockServer {
	var mockServer mock_server.MockServer
	mockServer.Start()
	fmt.Println("Mock server is running at " + mockServer.URL)
	return &mockServer
}

func startProxyServer() *ProxyServer {
	var proxyServer ProxyServer
	proxyServer.Setup()
	go proxyServer.Run()
	return &proxyServer
}

func TestServerReceivesBodyFromClient(t *testing.T) {
	// setup

	mockServer := startMockServer()
	defer mockServer.Close()

	proxyServer := startProxyServer()
	defer proxyServer.Close()

	clientRequestSender := ClientRequestSender{
		ServerURL: mockServer.URL,
		ProxyURL: proxyServer.URL,
	}

	// test
	clientRequestSender.SendRequestWithBody("TestServerReceivesBodyFromClient")
	clientRequestSender.WaitForAllRequests()

	if mockServer.LastRequest.BodyText != "TestServerReceivesBodyFromClient" {
		t.Error("Unexpected request body at Server: " + mockServer.LastRequest.BodyText)
	}
}

func TestClientReceivesTheBodyFromServer(t *testing.T) {
	// setup
	mockServer := startMockServer()
	defer mockServer.Close()

	proxyServer := startProxyServer()
	defer proxyServer.Close()

	clientRequestSender := ClientRequestSender{
		ServerURL: mockServer.URL,
		ProxyURL: proxyServer.URL,
	}

	// test
	clientRequestSender.SendRequest(func(_ *http.Response, responseBodyText string) {
		// Compare ALL fields of Response, including headers and the code!
		if responseBodyText != "ololo-shmololo\n" {
			t.Error("Unexpected response body:" + string(responseBodyText))
		}
	})

	clientRequestSender.WaitForAllRequests()
}
