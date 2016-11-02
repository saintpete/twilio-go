package twilio

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestTranscriptionDelete(t *testing.T) {
	t.Parallel()
	called := false
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(204)
	}))
	defer s.Close()
	client := NewClient("AC123", "123", nil)
	client.Base = s.URL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Transcriptions.Delete(ctx, "TR4c7f9a71f19b7509cb1e8344af78fc82")
	if err != nil {
		t.Fatal(err)
	}
	if called == false {
		t.Error("never hit server")
	}
}

func TestTranscriptionDeleteTwice(t *testing.T) {
	t.Parallel()
	client, server := getServerCode(transcriptionDeletedTwice, 404)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Transcriptions.Delete(ctx, "TR4c7f9a71f19b7509cb1e8344af78fc82")
	if err != nil {
		t.Fatal(err)
	}
}
