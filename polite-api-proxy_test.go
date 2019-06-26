package main

import (
	"testing"
	"fmt"
	"net/http"
	// "time"
	"strings"
	"log"
	"os"

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

func TestServerReceivesAllRequestSectionsFromClient(t *testing.T) {
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
	requestSentFromClient, error := http.NewRequest(
		"POST",
		mockServer.URL + "/some_path?query_key1=query_value1&query_key2=query_value2",
		strings.NewReader("TestServerReceivesAllRequestSectionsFromClientBody"),
	)

	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}

	requestSentFromClient.Header.Add("X-Test", "TestServerReceivesAllRequestSectionsFromClientHeader")

	clientRequestSender.SendRequest(requestSentFromClient, nil)
	clientRequestSender.WaitForAllRequests()

	requestReceivedByServer := mockServer.LastRequest.RawRequest

	if requestReceivedByServer.Method != requestSentFromClient.Method {
		t.Error("Methods do not match")
	}

	clientRequestURLWithoutHost := *requestSentFromClient.URL
	clientRequestURLWithoutHost.Host = ""
	clientRequestURLWithoutHost.Scheme = ""

	if *requestReceivedByServer.URL != clientRequestURLWithoutHost {
		t.Error("URLs do not match")
	}

	if requestReceivedByServer.Header.Get("X-Test") != "TestServerReceivesAllRequestSectionsFromClientHeader" {
		t.Error("Headers don't match")
	}

	if mockServer.LastRequest.BodyText != "TestServerReceivesAllRequestSectionsFromClientBody" {
		t.Error("Bodies don't match")
	}
}

func TestClientReceivesAllResponseSectionsFromServer(t *testing.T) {
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
	clientRequestSender.SendSimplestRequest(func(_ *http.Response, responseBodyText string) {
		// Compare ALL fields of Response, including headers and the code!
		if responseBodyText != "ololo-shmololo\n" {
			t.Error("Unexpected response body:" + string(responseBodyText))
		}
	})

	clientRequestSender.WaitForAllRequests()
}
