package twilio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var envClient = NewClient(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), nil)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "SM26b3b00f8def53be77c5697183bfe95e"
	msg, err := envClient.Messages.Get(sid)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Sid != sid {
		t.Errorf("expected Sid to equal %s, got %s", sid, msg.Sid)
	}
}

func TestGetPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	page, err := envClient.Messages.GetPage(url.Values{"PageSize": []string{"5"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Messages) != 5 {
		t.Fatalf("expected len(messages) to be 5, got %d", len(page.Messages))
	}
}

const from = "+19253920364"
const to = "+19253920364"

func TestSendMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	msg, err := envClient.Messages.SendMessage(from, to, "twilio-go testing!", nil)
	if err != nil {
		t.Fatal(err)
	}
	if msg.From != from {
		t.Errorf("expected From to be from, got error")
	}
}

func TestGetMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	msg, err := envClient.Messages.Get("SM5a52bc49b2354703bfdea7e92b44b385")
	if err != nil {
		t.Fatal(err)
	}
	if msg.ErrorCode != CodeUnknownDestination {
		t.Errorf("expected Code to be %d, got %d", CodeUnknownDestination, msg.ErrorCode)
	}
	if msg.ErrorMessage == "" {
		t.Errorf(`expected ErrorMessage to be non-empty, got ""`)
	}
}

func TestIterateAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	iter := envClient.Messages.GetPageIterator(url.Values{"PageSize": []string{"500"}})
	count := 0
	for {
		_, err := iter.Next()
		if err == NoMoreResults {
			break
		}
		if err != nil {
			t.Fatal(err)
			break
		}
		count++
		if count > 15 {
			fmt.Println("count > 15")
			t.Fail()
			break
		}
	}
}

func TestGetMediaURLs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	urls, err := envClient.Messages.GetMediaURLs("MM89a8c4a6891c53054e9cd604922bfb61", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Errorf("Wrong number of URLs returned: %d", len(urls))
	}
	if !strings.HasPrefix(urls[0].String(), "https://s3-external-1.amazonaws.com/media.twiliocdn.com/"+sid) {
		t.Errorf("wrong url: %s", urls[0].String())
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()
	str := `{"sid": "SM26b3b00f8def53be77c5697183bfe95e", "date_created": "Tue, 20 Sep 2016 22:59:57 +0000", "date_updated": "Tue, 20 Sep 2016 22:59:57 +0000", "date_sent": "Tue, 20 Sep 2016 22:59:57 +0000", "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279", "to": "+13365584092", "from": "+19253920364", "messaging_service_sid": null, "body": "Welcome to ZomboCom.", "status": "delivered", "num_segments": "1", "num_media": "0", "direction": "outbound-reply", "api_version": "2010-04-01", "price": "-0.00750", "price_unit": "USD", "error_code": null, "error_message": null, "uri": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Messages/SM26b3b00f8def53be77c5697183bfe95e.json", "subresource_uris": {"media": "/2010-04-01/Accounts/AC58f1e8f2b1c6b88ca90a012a4be0c279/Messages/SM26b3b00f8def53be77c5697183bfe95e/Media.json"}}`
	msg := new(Message)
	if err := json.Unmarshal([]byte(str), &msg); err != nil {
		t.Fatal(err)
	}
	if msg.Sid != "SM26b3b00f8def53be77c5697183bfe95e" {
		t.Errorf("wrong sid")
	}
	got := msg.DateCreated.Time.Format(time.RFC3339)
	want := "2016-09-20T22:59:57Z"
	if got != want {
		t.Errorf("msg.DateCreated: got %s, want %s", got, want)
	}
	if msg.Direction != DirectionOutboundReply {
		t.Errorf("wrong direction")
	}
	if msg.Status != StatusDelivered {
		t.Errorf("wrong status")
	}
	if msg.Body != "Welcome to ZomboCom." {
		t.Errorf("wrong body")
	}
	if msg.From != PhoneNumber("+19253920364") {
		t.Errorf("wrong from")
	}
	if msg.FriendlyPrice() != "$0.0075" {
		t.Errorf("wrong friendly price %v, want %v", msg.FriendlyPrice(), "$0.00750")
	}
}

func TestStatusFriendly(t *testing.T) {
	t.Parallel()
	if StatusQueued.Friendly() != "Queued" {
		t.Errorf("expected StatusQueued.Friendly to equal Queued, got %s", StatusQueued.Friendly())
	}
	s := Status("in-progress")
	if f := s.Friendly(); f != "In Progress" {
		t.Errorf("expected In Progress.Friendly to equal In Progress, got %s", f)
	}
}
