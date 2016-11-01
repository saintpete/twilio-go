package twilio

import (
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestGetCallerID(t *testing.T) {
	t.Parallel()
	client, server := getServer(callerIDInstance)
	defer server.Close()
	id, err := client.OutgoingCallerIDs.Get(context.Background(), "PNca86cf94c7d4f89e0bd45bfa7d9b9e7d")
	if err != nil {
		t.Fatal(err)
	}
	if id.PhoneNumber.Friendly() != "+1 925-271-7005" {
		t.Errorf("got bad phone number, want %s", id.PhoneNumber.Friendly())
	}
}

func TestVerifyCallerID(t *testing.T) {
	t.Parallel()
	client, server := getServer(callerIDVerify)
	defer server.Close()
	data := url.Values{}
	data.Set("PhoneNumber", "+14105551234")
	data.Set("FriendlyName", "test friendly name")
	id, err := client.OutgoingCallerIDs.Create(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if id.PhoneNumber.Friendly() != "+1 410-555-1234" {
		t.Errorf("got bad phone number, want %s", id.PhoneNumber.Friendly())
	}
	if len(id.ValidationCode) != 6 {
		t.Errorf("expected 6 digit validation code, got %s", id.ValidationCode)
	}
}
