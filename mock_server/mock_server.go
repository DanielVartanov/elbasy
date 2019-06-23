package mock_server

import (
	"net/http/httptest"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
)

type DebuggableRequest struct {
	BodyText string
	RawRequest *http.Request
}

type MockServer struct {
	URL string
	LastRequest DebuggableRequest

	httptestServer *httptest.Server
}

func (mockServer *MockServer) Start() {
	mockServer.httptestServer = httptest.NewServer(
		http.HandlerFunc(
			func (responseWriter http.ResponseWriter, request *http.Request) {
				fmt.Println("Received a request at Mock server")

				requestBodyText, error := ioutil.ReadAll(request.Body)
				if error != nil {
					fmt.Println(error)
					os.Exit(1)
				}
				mockServer.LastRequest = DebuggableRequest{RawRequest: request, BodyText: string(requestBodyText)}

				fmt.Printf("  RequestURI = %s, URL = %s, Method = %s, Body = \"%s\"\n", request.RequestURI, request.URL.String(), request.Method, requestBodyText)
				fmt.Printf("mockServer.LastRequest.BodyText = %s\n", mockServer.LastRequest.BodyText)
				fmt.Fprintln(responseWriter, "ololo-shmololo")
			},
		),
	)
	mockServer.URL = mockServer.httptestServer.URL
}

func (mockServer *MockServer) Close() {
	mockServer.httptestServer.Close()
}
