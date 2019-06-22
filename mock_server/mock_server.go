package mock_server

import (
	"net/http/httptest"
	"net/http"
	"fmt"
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
