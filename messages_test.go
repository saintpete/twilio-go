package twilio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestGetPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	c := NewClient(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), nil)
	page, err := c.Messages.GetPage(url.Values{"PageSize": []string{"5"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Messages) != 5 {
		t.Fatalf("expected len(messages) to be 5, got %d", len(page.Messages))
	}
}

func TestIterateAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	c := NewClient(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), nil)
	iter := c.Messages.GetPageIterator(url.Values{"PageSize": []string{"500"}})
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

func TestDecode(t *testing.T) {
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
}

func TestStatusFriendly(t *testing.T) {
	if StatusQueued.Friendly() != "Queued" {
		t.Errorf("expected StatusQueued.Friendly to equal Queued, got %s", StatusQueued.Friendly())
	}
}
