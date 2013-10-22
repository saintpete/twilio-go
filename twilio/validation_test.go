package twilio

import (
	"fmt"
	"net/url"
	"testing"
)

func TestClientValidateIncomingRequest(t *testing.T) {
	// Based on example at https://www.twilio.com/docs/security#validating-requests
	authToken := "12345"
	twilioClient := CreateClient("", authToken, nil)
	host := "https://mycompany.com"
	URL := "/myapp.php?foo=1&bar=2"
	xTwilioSignature := "RSOYDt4T1cUTdK1PDd93/VVr8B8="
	postForm := url.Values{
		"Digits":  {"1234"},
		"To":      {"+18005551212"},
		"From":    {"+14158675309"},
		"Caller":  {"+14158675309"},
		"CallSid": {"CA1234567890ABCDE"},
	}

	err := twilioClient.validateIncomingRequest(host, URL, postForm, xTwilioSignature)
	if err != nil {
		fmt.Println("Unexpected error:", err)
		t.Fail()
	}

	URL += "&cat=3"
	err = twilioClient.validateIncomingRequest(host, URL, postForm, xTwilioSignature)
	if err == nil {
		fmt.Println("Expected an error but got none")
		t.Fail()
	}
}
