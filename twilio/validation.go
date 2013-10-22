package twilio

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sort"
)

// See https://www.twilio.com/docs/security#validating-requests for more information
func (c *Client) ValidateIncomingRequest(req *http.Request) (err error) {
	err = req.ParseForm()
	if err != nil {
		return
	}

	return c.validateIncomingRequest(req.URL.String(), req.Form, req.Header.Get("X-Twilio-Signature"))
}

func (c *Client) validateIncomingRequest(URL string, postForm url.Values, xTwilioSignature string) (err error) {
	// Take the full URL of the request URL you specify for your
	// phone number or app, from the protocol (https...) through
	// the end of the query string (everything after the ?).
	str := URL

	// If the request is a POST, sort all of the POST parameters
	// alphabetically (using Unix-style case-sensitive sorting order).
	keys := make([]string, 0, len(postForm))
	for key := range postForm {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Iterate through the sorted list of POST parameters, and append
	// the variable name and value (with no delimiters) to the end
	// of the URL string.
	for _, key := range keys {
		str += key + postForm[key][0]
	}

	// Sign the resulting string with HMAC-SHA1 using your AuthToken
	// as the key (remember, your AuthToken's case matters!).
	mac := hmac.New(sha1.New, []byte(c.AuthToken))
	mac.Write([]byte(str))
	expectedMac := mac.Sum(nil)

	// Base64 encode the resulting hash value.
	expectedSignature := base64.StdEncoding.EncodeToString(expectedMac)

	// Compare your hash to ours, submitted in the X-Twilio-Signature header.
	// If they match, then you're good to go.
	if xTwilioSignature != expectedSignature {
		return errors.New("Bad X-Twilio-Signature")
	}

	return nil
}
