package main

import (
	"net/url"
	"net/http"
	"sync"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"strings"
)

type ClientRequestSender struct {
	ServerURL string
	ProxyURL string

	requestsWaitGroup sync.WaitGroup
}

func (self *ClientRequestSender) composeSimplestRequest() *http.Request {
	request, error := http.NewRequest("GET", self.ServerURL, nil)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
	return request
}

func (self *ClientRequestSender) composeRequestWithBody(requestBodyText string) *http.Request {
	request, error := http.NewRequest(
		"POST",
		self.ServerURL,
		strings.NewReader(requestBodyText),
	)

	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}
	return request
}

func (self *ClientRequestSender) SendRequest(callback func(response *http.Response, responseBodyText string)) {
	self.requestsWaitGroup.Add(1)
	go self.actuallySendRequest(self.composeSimplestRequest(), callback)
}

func (self *ClientRequestSender) SendRequestWithBody(requestBodyText string) { // Just accept http.Request
	self.requestsWaitGroup.Add(1)
	go self.actuallySendRequest(self.composeRequestWithBody(requestBodyText), nil)
}

func (self *ClientRequestSender) buildHTTPClient() http.Client {
	proxyUrl, error := url.Parse(self.ProxyURL)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
		Proxy: http.ProxyURL(proxyUrl),
	}

	return http.Client{Transport: transport}
}

func (self *ClientRequestSender) actuallySendRequest(request *http.Request, callback func(response *http.Response, responseBodyText string)) {
	fmt.Println("Sending a client request to the server via the proxy...")

	client := self.buildHTTPClient()

	response, error := client.Do(request)
	if error != nil {
		log.Fatal(error)
		os.Exit(1)
	}

	defer response.Body.Close()
	responseBodyText, error := ioutil.ReadAll(response.Body)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	if callback != nil {
		callback(response, string(responseBodyText))
	}

	self.requestsWaitGroup.Done()
}

func (self *ClientRequestSender) WaitForAllRequests() {
	self.requestsWaitGroup.Wait()
}
