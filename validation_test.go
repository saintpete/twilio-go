package twilio

import (
	"net/url"
	"testing"
)

func TestClientValidateIncomingRequest(t *testing.T) {
	t.Parallel()
	// Based on example at https://www.twilio.com/docs/security#validating-requests
	authToken := "12345"
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

	err := validateIncomingRequest(host, authToken, URL, postForm, xTwilioSignature)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	URL += "&cat=3"
	err = validateIncomingRequest(host, authToken, URL, postForm, xTwilioSignature)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}
