package main

import (
	"net/url"
	"net/http"
	"sync"
	"io/ioutil"
	"os"
	"fmt"
)

type ClientRequestSender struct {
	ServerURL string
	ProxyURL string

	requestsWaitGroup sync.WaitGroup
	setupOnceLock sync.Once
}

func (self *ClientRequestSender) setupProxyUsageOnce() {
	self.setupOnceLock.Do(func() {
		os.Setenv("HTTP_PROXY", self.ProxyURL) // Why env var not working?
		proxyUrl, error := url.Parse(self.ProxyURL)
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	})
}

func (self *ClientRequestSender) SendRequest(callback func(response *http.Response, responseBodyText string)) {
	self.setupProxyUsageOnce()
	self.requestsWaitGroup.Add(1)
	go self.actuallySendRequest(callback)
}

func (self *ClientRequestSender) WaitForAllRequests() {
	self.requestsWaitGroup.Wait()
}

func (self *ClientRequestSender) actuallySendRequest(callback func(response *http.Response, responseBodyText string)) {
	fmt.Println("Sending a client request to the server via the proxy...")
	response, error := http.Get(self.ServerURL)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	defer response.Body.Close()
	responseBodyText, error := ioutil.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	callback(response, string(responseBodyText))

	self.requestsWaitGroup.Done()
}
