package twilio

import (
	"net/http"
	"net/http/httptest"
	"os"
)

// the envClient is configured to use an Account Sid and Auth Token set in the
// environment. all non-short tests should use the envClient
var envClient = NewClient(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), nil)

// getServer returns a http server that returns the given bytes when requested,
// and a Client configured to make requests to that server.
func getServer(response []byte) (*Client, *httptest.Server) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := w.Write(response); err != nil {
			panic(err)
		}
	}))
	client := NewClient("AC123", "456", nil)
	client.Base = s.URL
	return client, s
}

var conferenceInstance = []byte(`{"sid": "CF169b5eebb07ec48e0f9f2ee904b385c5", "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279", "friendly_name": "testConference", "status": "completed", "date_created": "Fri, 23 Aug 2013 21:52:27 +0000", "api_version": "2010-04-01", "date_updated": "Fri, 23 Aug 2013 21:52:34 +0000", "region": "us1", "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5.json", "subresource_uris": {"participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5/Participants.json"}}`)
var conferenceInstanceSid = "CF169b5eebb07ec48e0f9f2ee904b385c5"

var conferencePage = []byte(`{"first_page_uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences.json?PageSize=3&Page=0", "end": 2, "conferences": [{"sid": "CF347aef00d0b0ba10eec6a77fabdd1c95", "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279", "friendly_name": "turbovote", "status": "completed", "date_created": "Fri, 13 Sep 2013 19:53:14 +0000", "api_version": "2010-04-01", "date_updated": "Fri, 13 Sep 2013 19:53:34 +0000", "region": "us1", "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF347aef00d0b0ba10eec6a77fabdd1c95.json", "subresource_uris": {"participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF347aef00d0b0ba10eec6a77fabdd1c95/Participants.json"}}, {"sid": "CF169b5eebb07ec48e0f9f2ee904b385c5", "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279", "friendly_name": "testConference", "status": "completed", "date_created": "Fri, 23 Aug 2013 21:52:27 +0000", "api_version": "2010-04-01", "date_updated": "Fri, 23 Aug 2013 21:52:34 +0000", "region": "us1", "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5.json", "subresource_uris": {"participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5/Participants.json"}}, {"sid": "CFb2b77e00f9e97764746aff575a15bfbb", "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279", "friendly_name": "testConference", "status": "completed", "date_created": "Fri, 23 Aug 2013 21:48:44 +0000", "api_version": "2010-04-01", "date_updated": "Fri, 23 Aug 2013 21:52:01 +0000", "region": "us1", "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CFb2b77e00f9e97764746aff575a15bfbb.json", "subresource_uris": {"participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CFb2b77e00f9e97764746aff575a15bfbb/Participants.json"}}], "previous_page_uri": null, "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences.json?PageSize=3", "page_size": 3, "start": 0, "next_page_uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences.json?PageSize=3&Page=1&AfterSid=CFb2b77e00f9e97764746aff575a15bfbb", "page": 0}`)
