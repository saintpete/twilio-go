package twilio

import (
	"testing"

	"golang.org/x/net/context"
)

func TestApplicationGet(t *testing.T) {
	t.Parallel()
	client, server := getServer(applicationInstance)
	defer server.Close()
	application, err := client.Applications.Get(context.Background(), "AP7d6fd7b9a8894e36877dc2355da381c8")
	if err != nil {
		t.Fatal(err)
	}
	if application.FriendlyName != "Hackpack for Heroku and Flask" {
		t.Error("bad friendly name")
	}
}
