package main

// Idea: use httptest.ResponseRecorders
/*
import (
	"testing"
	"fmt"
	"net/http"
	"strings"
	"log"
	"os"

	"./mock_server"
)


func startMockServer(responseCompositionFunction func (responseWriter http.ResponseWriter, request *http.Request)) *mock_server.MockServer {
	var mockServer mock_server.MockServer
	mockServer.Start(responseCompositionFunction)
	fmt.Println("Mock server is running at " + mockServer.URL)
	return &mockServer
}

func startProxyServer() *ProxyServer {
	var proxyServer ProxyServer
	proxyServer.Setup()
	go proxyServer.Run()
	return &proxyServer
}

func TestServerReceivesAllRequestSectionsFromClient(t *testing.T) {
	// setup
	mockServer := startMockServer(nil)
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
	// setup
	mockServer := startMockServer(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Add("X-Test-Mock-Server", "TestClientReceivesAllResponseSectionsFromServerHeader")
		responseWriter.WriteHeader(201)
		fmt.Fprint(responseWriter, "TestClientReceivesAllResponseSectionsFromServerBody")
	})
	defer mockServer.Close()

	proxyServer := startProxyServer()
	defer proxyServer.Close()

	clientRequestSender := ClientRequestSender{
		ServerURL: mockServer.URL,
		ProxyURL: proxyServer.URL,
	}

	// test
	clientRequestSender.SendSimplestRequest(func(response *http.Response, responseBodyText string) {

		if response.Proto != "HTTP/1.1" {
			t.Error("Protocol versions do not match")
		}

		if response.StatusCode != 201 {
			t.Error("Status codes do not match")
		}

		if response.Header.Get("X-Test-Mock-Server") != "TestClientReceivesAllResponseSectionsFromServerHeader" {
			t.Error("Headers do not match")
		}

		if responseBodyText != "TestClientReceivesAllResponseSectionsFromServerBody" {
			t.Error("Bodies do not match")
		}


	})

	clientRequestSender.WaitForAllRequests()
}
*/
