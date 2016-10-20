package twilio

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kevinburke/rest"
)

// invalid status here on purpose to check we use a different one.
var notFoundResp = []byte("{\"code\": 20404, \"message\": \"The requested resource /2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Calls/unknown.json was not found\", \"more_info\": \"https://www.twilio.com/docs/errors/20404\", \"status\": 428}")

var errorServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(404)
	w.Write(notFoundResp)
}))

func Test404Error(t *testing.T) {
	client := NewClient("", "", nil)
	client.Base = errorServer.URL
	sid := "unknown"
	_, err := client.Calls.Get(sid)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	rerr, ok := err.(*rest.Error)
	if !ok {
		t.Fatalf("expected to convert err %v to rest.Error, couldn't", err)
	}
	if !strings.Contains(rerr.Title, "The requested resource /2010-04-01") {
		t.Errorf("expected Title to contain 'The requested resource', got %s", rerr.Title)
	}
	if rerr.ID != "20404" {
		t.Errorf("expected ID to be 20404, got %s", rerr.ID)
	}
	if rerr.Type != "https://www.twilio.com/docs/errors/20404" {
		t.Errorf("expected Type to be a Twilio URL, got %s", rerr.Type)
	}
	if rerr.StatusCode != 404 {
		t.Errorf("expected StatusCode to be 404, got %d", rerr.StatusCode)
	}
}
