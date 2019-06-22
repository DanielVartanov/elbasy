package main


import (
	"testing"
	"os"
	// "net"
	"net/http"
	"net/http/httptest"
	"fmt"
	"io/ioutil"
	"time"
	"net/url"
)

type MockServer struct {
	URL string

	httptestServer *httptest.Server
}

func (mockServer *MockServer) Start() {
	mockServer.httptestServer = httptest.NewServer(
		http.HandlerFunc(
			func (responseWriter http.ResponseWriter, request *http.Request) {
				fmt.Println("Received a request at Mock server")
				fmt.Fprintln(responseWriter, "ololo-shmololo")
			},
		),
	)
	mockServer.URL = mockServer.httptestServer.URL
}

func (mockServer *MockServer) Close() {
	mockServer.httptestServer.Close()
}

func TestProxyServer(t *testing.T) {
	var mockServer MockServer
	mockServer.Start()
	defer mockServer.Close()
	fmt.Println("Mock server is running at " + mockServer.URL)

	proxyServer := RunProxyServer()
	defer proxyServer.Close()

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
