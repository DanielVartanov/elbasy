package mockremotehost

import (
	"testing"
	"net/http"
)

var mrh *Server = NewServer("mockremotehost.lvh.me", 9876, "mockremotehost.lvh.me.pem", "mockremotehost.lvh.me-key.pem")

func TestURL(t *testing.T) {
	if mrh.URL().String() != "https://mockremotehost.lvh.me:9876" {
		t.Fail()
	}
}

func TestOnRequest(t *testing.T) {
	err := mrh.BindToPort()
	if err != nil { t.Fatalf("mrh.BindToPort(): %v", err) }

	go func() {
		err = mrh.AcceptConnections()
		if err != nil { t.Logf("mrh.AcceptConnections(): %v", err) }
	}()

	mrh.OnRequest(func(w http.ResponseWriter, req *http.Request){
		if req.URL.Path != "/expected_path" {
			t.Fatalf("`request` param in OnRequest does not equal to actual request")
		}

		w.Header().Set("X-Expected-Key", "expected value")
		w.WriteHeader(http.StatusOK)
	})

	response, err := http.Get(mrh.URL().String() + "/expected_path")
	if err != nil { t.Fatalf("http.Get(): %v", err) }

	if response.Header.Get("X-Expected-Key") != "expected value" {
		t.Fatalf("Server did not respond with the response written in OnRequest")
	}

	err = mrh.Stop()
	if err != nil { t.Fatalf("mrh.Close(): %v", err) }
}
