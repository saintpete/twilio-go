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

// useful trick: highlight the JSON range and hit `python -m json.tool` to
// pretty format it.

var conferenceInstance = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "api_version": "2010-04-01",
    "date_created": "Fri, 23 Aug 2013 21:52:27 +0000",
    "date_updated": "Fri, 23 Aug 2013 21:52:34 +0000",
    "friendly_name": "testConference",
    "region": "us1",
    "sid": "CF169b5eebb07ec48e0f9f2ee904b385c5",
    "status": "completed",
    "subresource_uris": {
        "participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5/Participants.json"
    },
    "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5.json"
}
`)
var conferenceInstanceSid = "CF169b5eebb07ec48e0f9f2ee904b385c5"

var conferencePage = []byte(`
{
    "conferences": [
        {
            "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
            "api_version": "2010-04-01",
            "date_created": "Fri, 13 Sep 2013 19:53:14 +0000",
            "date_updated": "Fri, 13 Sep 2013 19:53:34 +0000",
            "friendly_name": "turbovote",
            "region": "us1",
            "sid": "CF347aef00d0b0ba10eec6a77fabdd1c95",
            "status": "completed",
            "subresource_uris": {
                "participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF347aef00d0b0ba10eec6a77fabdd1c95/Participants.json"
            },
            "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF347aef00d0b0ba10eec6a77fabdd1c95.json"
        },
        {
            "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
            "api_version": "2010-04-01",
            "date_created": "Fri, 23 Aug 2013 21:52:27 +0000",
            "date_updated": "Fri, 23 Aug 2013 21:52:34 +0000",
            "friendly_name": "testConference",
            "region": "us1",
            "sid": "CF169b5eebb07ec48e0f9f2ee904b385c5",
            "status": "completed",
            "subresource_uris": {
                "participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5/Participants.json"
            },
            "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CF169b5eebb07ec48e0f9f2ee904b385c5.json"
        },
        {
            "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
            "api_version": "2010-04-01",
            "date_created": "Fri, 23 Aug 2013 21:48:44 +0000",
            "date_updated": "Fri, 23 Aug 2013 21:52:01 +0000",
            "friendly_name": "testConference",
            "region": "us1",
            "sid": "CFb2b77e00f9e97764746aff575a15bfbb",
            "status": "completed",
            "subresource_uris": {
                "participants": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CFb2b77e00f9e97764746aff575a15bfbb/Participants.json"
            },
            "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences/CFb2b77e00f9e97764746aff575a15bfbb.json"
        }
    ],
    "end": 2,
    "first_page_uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences.json?PageSize=3&Page=0",
    "next_page_uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences.json?PageSize=3&Page=1&AfterSid=CFb2b77e00f9e97764746aff575a15bfbb",
    "page": 0,
    "page_size": 3,
    "previous_page_uri": null,
    "start": 0,
    "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Conferences.json?PageSize=3"
}
`)

var sendMessageResponse = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "api_version": "2010-04-01",
    "body": "twilio-go testing!",
    "date_created": "Wed, 26 Oct 2016 18:12:20 +0000",
    "date_sent": null,
    "date_updated": "Wed, 26 Oct 2016 18:12:20 +0000",
    "direction": "outbound-api",
    "error_code": null,
    "error_message": null,
    "from": "+19253920364",
    "messaging_service_sid": null,
    "num_media": "0",
    "num_segments": "1",
    "price": null,
    "price_unit": "USD",
    "sid": "SM9b7db369463c439384fe9abb09751af8",
    "status": "queued",
    "subresource_uris": {
        "media": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Messages/SM9b7db369463c439384fe9abb09751af8/Media.json"
    },
    "to": "+19253920364",
    "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Messages/SM9b7db369463c439384fe9abb09751af8.json"
}
`)

const from = "+19253920364"
const to = "+19253920364"
