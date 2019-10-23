package mockremotehost

import (
	"testing"
	"net/http"
)

const port = 9000

func TestURL(t *testing.T) {
	mrh := Server{Port: 9876}
	if mrh.URL().String() != "http://localhost:9876" {
		t.Fail()
	}
}

func TestOnRequest(t *testing.T) {
	var mrh Server = Server{Port: port}
	err := mrh.BindToPort()
	if err != nil { t.Errorf("mrh.BindToPort(): %v", err) }

	go func() {
		err = mrh.AcceptConnections()
		if err != nil { t.Errorf("mrh.AcceptConnections(): %v", err) }
	}()

	mrh.OnRequest(func(w http.ResponseWriter, req *http.Request){
		if req.URL.Path != "/expected_path" {
			t.Errorf("`request` param in OnRequest does not equal to actual request")
		}

		w.Header().Set("X-Expected-Key", "expected value")
		w.WriteHeader(http.StatusOK)
	})

	response, err := http.Get(mrh.URL().String() + "/expected_path")
	if err != nil { t.Errorf("http.Get(): %v", err) }

	if response.Header.Get("X-Expected-Key") != "expected value" {
		t.Errorf("Server did not respond with the response written in OnRequest")
	}
}
