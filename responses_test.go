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
	client.Monitor.Base = s.URL
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

var makeCallResponse = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "annotation": null,
    "answered_by": null,
    "api_version": "2010-04-01",
    "caller_name": null,
    "date_created": null,
    "date_updated": null,
    "direction": "outbound-api",
    "duration": null,
    "end_time": null,
    "forwarded_from": null,
    "from": "+19253920364",
    "from_formatted": "(925) 392-0364",
    "group_sid": null,
    "parent_call_sid": null,
    "phone_number_sid": "PN5fb9ed903e184c8baa86c1fb7544ca0f",
    "price": null,
    "price_unit": "USD",
    "sid": "CA47b862ce3b99a6d79939320a9aa54a02",
    "start_time": null,
    "status": "queued",
    "subresource_uris": {
        "notifications": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Calls/CA47b862ce3b99a6d79939320a9aa54a02/Notifications.json",
        "recordings": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Calls/CA47b862ce3b99a6d79939320a9aa54a02/Recordings.json"
    },
    "to": "+19252717005",
    "to_formatted": "(925) 271-7005",
    "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Calls/CA47b862ce3b99a6d79939320a9aa54a02.json"
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

var getMessageResponse = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "api_version": "2010-04-01",
    "body": "Welcome to ZomboCom.",
    "date_created": "Tue, 20 Sep 2016 22:59:57 +0000",
    "date_sent": "Tue, 20 Sep 2016 22:59:57 +0000",
    "date_updated": "Tue, 20 Sep 2016 22:59:57 +0000",
    "direction": "outbound-reply",
    "error_code": null,
    "error_message": null,
    "from": "+19253920364",
    "messaging_service_sid": null,
    "num_media": "0",
    "num_segments": "1",
    "price": "-0.00750",
    "price_unit": "USD",
    "sid": "SM26b3b00f8def53be77c5697183bfe95e",
    "status": "delivered",
    "subresource_uris": {
        "media": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Messages/SM26b3b00f8def53be77c5697183bfe95e/Media.json"
    },
    "to": "+13365584092",
    "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Messages/SM26b3b00f8def53be77c5697183bfe95e.json"
}
`)

var alertInstanceResponse = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "alert_text": "Msg&sourceComponent=12000&ErrorCode=11200&httpResponse=405&url=https%3A%2F%2Fkev.inburke.com%2Fzombo%2Fzombocom.mp3&LogLevel=ERROR",
    "api_version": "2010-04-01",
    "date_created": "2016-10-27T02:34:21Z",
    "date_generated": "2016-10-27T02:34:21Z",
    "date_updated": "2016-10-27T02:34:23Z",
    "error_code": "11200",
    "log_level": "error",
    "more_info": "https://www.twilio.com/docs/errors/11200",
    "request_headers": null,
    "request_method": "POST",
    "request_url": "https://kev.inburke.com/zombo/zombocom.mp3",
    "request_variables": "Called=%2B19252717005&ToState=CA&CallerCountry=US&Direction=outbound-api&CallerState=CA&ToZip=94596&CallSid=CA6d27370cbbfb605521fe8800bb73f2d2&To=%2B19252717005&CallerZip=94514&ToCountry=US&ApiVersion=2010-04-01&CalledZip=94596&CalledCity=PLEASANTON&CallStatus=in-progress&From=%2B19253920364&AccountSid=AC58f1e8f2b1c6b88ca90a012a4be0c279&CalledCountry=US&CallerCity=BRENTWOOD&Caller=%2B19253920364&FromCountry=US&ToCity=PLEASANTON&FromCity=BRENTWOOD&CalledState=CA&FromZip=94514&FromState=CA",
    "resource_sid": "CA6d27370cbbfb605521fe8800bb73f2d2",
    "response_body": "<html>\r\n<head><title>405 Not Allowed</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>405 Not Allowed</h1></center>\r\n<hr><center>nginx</center>\r\n</body>\r\n</html>",
    "response_headers": "Transfer-Encoding=chunked&Server=cloudflare-nginx&CF-RAY=2f82bf9cb8102204-EWR&Set-Cookie=__cfduid%3Dd46f1cfd57d664c3038ae66f1c1de9e751477535661%3B+expires%3DFri%2C+27-Oct-17+02%3A34%3A21+GMT%3B+path%3D%2F%3B+domain%3D.inburke.com%3B+HttpOnly&Date=Thu%2C+27+Oct+2016+02%3A34%3A21+GMT&Content-Type=text%2Fhtml",
    "service_sid": null,
    "sid": "NO00ed1fb4aa449be2434d54ec8e492349",
    "url": "https://monitor.twilio.com/v1/Alerts/NO00ed1fb4aa449be2434d54ec8e492349"
}
`)

var alertListResponse = []byte(`
{
    "alerts": [
        {
            "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
            "alert_text": "Msg&sourceComponent=12000&ErrorCode=11200&httpResponse=405&url=https%3A%2F%2Fkev.inburke.com%2Fzombo%2Fzombocom.mp3&LogLevel=ERROR",
            "api_version": "2010-04-01",
            "date_created": "2016-10-27T02:34:21Z",
            "date_generated": "2016-10-27T02:34:21Z",
            "date_updated": "2016-10-27T02:34:23Z",
            "error_code": "11200",
            "log_level": "error",
            "more_info": "https://www.twilio.com/docs/errors/11200",
            "request_method": "POST",
            "request_url": "https://kev.inburke.com/zombo/zombocom.mp3",
            "resource_sid": "CA6d27370cbbfb605521fe8800bb73f2d2",
            "service_sid": null,
            "sid": "NO00ed1fb4aa449be2434d54ec8e492349",
            "url": "https://monitor.twilio.com/v1/Alerts/NO00ed1fb4aa449be2434d54ec8e492349"
        },
        {
            "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
            "alert_text": "sourceComponent=14100&ErrorCode=14101&LogLevel=ERROR&Msg=The+destination+number+for+a+TwiML+message+can+not+be+the+same+as+the+originating+number+of+an+incoming+message.&EmailNotification=false",
            "api_version": "2008-08-01",
            "date_created": "2016-10-26T18:12:20Z",
            "date_generated": "2016-10-26T18:12:20Z",
            "date_updated": "2016-10-26T18:12:24Z",
            "error_code": "14101",
            "log_level": "error",
            "more_info": "https://www.twilio.com/docs/errors/14101",
            "request_method": "POST",
            "request_url": "https://kev.inburke.com/zombo/zombo.php",
            "resource_sid": "SM07bb932105b6c969574be9ca5771d139",
            "service_sid": null,
            "sid": "NOffde6702f5035789b395f8618f0aa65a",
            "url": "https://monitor.twilio.com/v1/Alerts/NOffde6702f5035789b395f8618f0aa65a"
        }
    ],
    "meta": {
        "first_page_url": "https://monitor.twilio.com/v1/Alerts?PageSize=2&Page=0",
        "key": "alerts",
        "next_page_url": "https://monitor.twilio.com/v1/Alerts?PageSize=2&Page=1&PageToken=PANOffde6702f5035789b395f8618f0aa65a",
        "page": 0,
        "page_size": 2,
        "previous_page_url": null,
        "url": "https://monitor.twilio.com/v1/Alerts?PageSize=2&Page=0"
    }
}
`)

const from = "+19253920364"
const to = "+19253920364"
