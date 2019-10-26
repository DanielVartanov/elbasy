package main

import (
	"testing"
	"net"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"

	"github.com/DanielVartanov/elbasy/mockremotehost"
)

var mockserver *mockremotehost.Server
var elbasy *proxy
var client http.Client

func setup(t *testing.T) {
	// setup mock remote host
	mockserver = mockremotehost.NewServer("non-throttling.lvh.me", 9001, "non-throttling.lvh.me.pem", "non-throttling.lvh.me-key.pem")

	err := mockserver.BindToPort()
	if err != nil { t.Fatalf("mockserver.BindToPort(): %v", err) }

	go func() {
		err = mockserver.AcceptConnections()
		if err != nil { t.Logf("mockserver.AcceptConnections(): %v", err) }
	}()

	// setup elbasy server
	elbasy = newProxy("elbasy.lvh.me", 8443, "elbasy.lvh.me.pem", "elbasy.lvh.me-key.pem")

	err = elbasy.bindToPort()
	if err != nil { t.Fatalf("elbasy.bindToPort(): %v", err) }

	go func() {
		err = elbasy.acceptConnections()
		if err != nil { t.Logf("elbasy.acceptConnections(): %v", err) }
	}()

	// setup HTTP client
	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives: false,
		MaxIdleConnsPerHost: 100,
		Proxy: http.ProxyURL(elbasy.url()),
	}
	client = http.Client{Transport: transport}
}

func waitForSrvClosure(host string) {
	// Remove this method after the bug in golang stdlib is fixed: https://github.com/golang/go/issues/10527
	for {
		_, err := net.Dial("tcp", host)
		if err != nil { break }
	}
}

func teardown(t *testing.T) {
	err := mockserver.Stop()
	if err != nil { t.Logf("mockserver.Close(): %v", err) }

	err = elbasy.close()
	if err != nil { t.Logf("elbasy.close(): %v", err) }

	waitForSrvClosure(elbasy.url().Host)
	waitForSrvClosure(mockserver.URL().Host)
}

func sendRequest(t *testing.T, req *http.Request) *http.Response {
	resp, err := client.Do(req)
	if err != nil { t.Fatalf("client.Do(): %v", err) }
	return resp
}

func sendAnyRequest(t *testing.T) *http.Response {
	req, err := http.NewRequest("", mockserver.URL().String(), nil)
	if err != nil { t.Fatalf("http.NewRequest(): %v", err) }
	return sendRequest(t, req)
}

func TestRegularProxyRelaysExactRequest(t *testing.T) {
	setup(t)
	defer func() { teardown(t) }()

	sentReq, err := http.NewRequest(
		"POST",
		mockserver.URL().String() + "/some_path?somekey=somevalue",
		strings.NewReader("Some request body"),
	)
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	sentReq.Header.Set("X-Some-Header-Key", "Some header value")

	mockserver.OnRequest(func (w http.ResponseWriter, recvdReq *http.Request) {
		if recvdReq.Method != "POST" {
			t.Errorf("method doed not match: %v", recvdReq.Method)
		}

		if recvdReq.RequestURI != "/some_path?somekey=somevalue" {
			t.Errorf("path does not match: %v", recvdReq.RequestURI)
		}

		if recvdReq.Header.Get("X-Some-Header-Key") != "Some header value" {
			t.Errorf("header does not match: %v", recvdReq.Header.Get("X-Some-Header-Key"))
		}

		bytes, err := ioutil.ReadAll(recvdReq.Body)
		if err != nil { t.Fatalf("ioutil.ReadAll: %v", err) }
		body := string(bytes)
		if body != "Some request body" {
			t.Errorf("body does not match: %v", body)
		}
	})

	sendRequest(t, sentReq)
}

func TestRegularProxyRelaysExactResponse(t *testing.T) {
	setup(t)
	defer func() { teardown(t) }()

	mockserver.OnRequest(func (w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Some-Header-Key", "Some header value")
		w.WriteHeader(201)
		fmt.Fprint(w, "Some response body")
	})

	response := sendAnyRequest(t)

	if response.StatusCode != 201 {
		t.Errorf("status code does not match: %d", response.StatusCode)
	}

	if response.Header.Get("X-Some-Header-Key") != "Some header value" {
		t.Errorf("header does not match: %v", response.Header.Get("X-Some-Header-Key"))
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil { t.Fatalf("ioutil.ReadAll: %v", err) }
	body := string(bytes)
	if body != "Some response body" {
		t.Errorf("body does not match: %v", body)
	}
}
