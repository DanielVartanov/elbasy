package elbasy_test

import (
	"net/http"
	"net/http/httptest"
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

func (mockServer *MockServer) Start(responseCompositionFunction func (responseWriter http.ResponseWriter, request *http.Request)) {
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

				if responseCompositionFunction != nil {
					responseCompositionFunction(responseWriter, request)
				}
			},
		),
	)
	mockServer.URL = mockServer.httptestServer.URL
}

func (mockServer *MockServer) Close() {
	fmt.Println("Stopping a Mock server")
	mockServer.httptestServer.Close()
}
