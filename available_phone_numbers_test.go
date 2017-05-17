package twilio

import (
	"golang.org/x/net/context"
	"net/url"
	"testing"
)

func TestSearchAvailablePhoneNumbers(t *testing.T) {
	t.Parallel()
	client, server := getServer(availablePhoneNumbers)
	defer server.Close()

	data := url.Values{
		"Contains":     []string{"571*******"},
		"InRegion":     []string{"VA"},
		"SmsEnabled":   []string{"true"},
		"VoiceEnabled": []string{"true"},
	}

	res, err := client.AvailablePhoneNumbers.Local.GetPage(context.Background(), "US", data)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Numbers) != 1 {
		t.Errorf("expected 1 number, got %d", len(res.Numbers))
	}

	if res.Numbers[0].FriendlyName != "(510) 564-7903" {
		t.Errorf("unexpected friendly name: %s", res.Numbers[0].FriendlyName)
	}

	if res.Numbers[0].PhoneNumber != "+15105647903" {
		t.Errorf("unexpected phone number: %s", res.Numbers[0].PhoneNumber)
	}
}
