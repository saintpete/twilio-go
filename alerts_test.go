package twilio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGetAlert(t *testing.T) {
	t.Parallel()
	client, s := getServer(alertInstanceResponse)
	defer s.Close()
	sid := "NO00ed1fb4aa449be2434d54ec8e492349"
	alert, err := client.Monitor.Alerts.Get(context.Background(), sid)
	if err != nil {
		t.Fatal(err)
	}
	if alert.Sid != sid {
		t.Errorf("expected Sid to be %s, got %s", sid, alert.Sid)
	}
	if city := alert.RequestVariables.Get("CallerCity"); city != "BRENTWOOD" {
		t.Errorf("expected to get BRENTWOOD for CallerCity, got %s", city)
	}
}

func TestGetAlertPage(t *testing.T) {
	t.Parallel()
	client, s := getServer(alertListResponse)
	defer s.Close()
	data := url.Values{"PageSize": []string{"2"}}
	page, err := client.Monitor.Alerts.GetPage(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Alerts) != 2 {
		t.Errorf("expected 2 alerts, got %d", len(page.Alerts))
	}
	if page.Meta.Key != "alerts" {
		t.Errorf("expected Key to be 'alerts', got %s", page.Meta.Key)
	}
	if page.Meta.PageSize != 2 {
		t.Errorf("expected PageSize to be 2, got %d", page.Meta.PageSize)
	}
	if page.Meta.Page != 0 {
		t.Errorf("expected Page to be 0, got %d", page.Meta.Page)
	}
	if page.Meta.PreviousPageURL.Valid != false {
		t.Errorf("expected previousPage.Valid to be false, got true")
	}
	if page.Alerts[0].LogLevel != LogLevelError {
		t.Errorf("expected LogLevel to be Error, got %s", page.Alerts[0].LogLevel)
	}
	if page.Alerts[0].RequestMethod != "POST" {
		t.Errorf("expected RequestMethod to be 'POST', got %s", page.Alerts[0].RequestMethod)
	}
	if page.Alerts[0].ErrorCode != CodeHTTPRetrievalFailure {
		t.Errorf("expected ErrorCode to be '11200', got %d", page.Alerts[0].ErrorCode)
	}
}

func TestAlertFullPath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	// sigh, we need to hit the real URL for the Base stuff to work.
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	data := url.Values{"PageSize": []string{"2"}}
	iter := envClient.Monitor.Alerts.GetPageIterator(data)
	_, err := iter.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = iter.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// TODO figure out a good way to assert what URL gets hit; add a TestClient
	// or something.
}

func TestGetAlertIterator(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	// TODO make a better assertion here or a server that can return different
	// responses
	iter := envClient.Monitor.Alerts.GetPageIterator(url.Values{"PageSize": []string{"35"}})
	count := uint(0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for {
		page, err := iter.Next(ctx)
		if err == NoMoreResults {
			break
		}
		if err != nil {
			t.Fatal(err)
			break
		}
		if page.Meta.Page < count {
			t.Fatal("too small page number")
			return
		}
		count++
		if count > 15 {
			fmt.Println("count > 15")
			t.Fail()
			break
		}
	}
	if count < 10 {
		t.Errorf("Too small of a count - expected at least 10, got %d", count)
	}
}

func TestCapitalize(t *testing.T) {
	if capitalize("S") != "S" {
		t.Errorf("wrong")
	}
	if capitalize("s") != "S" {
		t.Errorf("wrong")
	}
	if capitalize("booo") != "Booo" {
		t.Errorf("wrong")
	}
}

var descriptionTests = []struct {
	in       []byte
	expected string
}{
	{alertDestination, "The destination number for a TwiML message can not be the same as the originating number of an incoming message"},
	{alert11200, "HTTP retrieval failure: status code 405 when fetching TwiML"},
	{alert14107, "Reply rate limit hit replying to +14156305833 from +19253920364"},
	{alert13225, "Forbidden phone number +886225475050"},
	{alert13227, "Not authorized to call +864008895080"},
	{alertUnknown, "Error 235342434: https://www.twilio.com/docs/errors/93455"},
}

func TestAlertDescription(t *testing.T) {
	for _, tt := range descriptionTests {
		alert := new(Alert)
		if err := json.Unmarshal(tt.in, alert); err != nil {
			panic(err)
		}
		if desc := alert.Description(); desc != tt.expected {
			t.Errorf("bad description: got %s, want %s", desc, tt.expected)
		}
	}
}

func TestAlertStatusCode(t *testing.T) {
	alert := new(Alert)
	alert.AlertText = "Msg&sourceComponent=12000&ErrorCode=11200&httpResponse=405&url=https%3A%2F%2Fkev.inburke.com%2Fzombo%2Fzombocom.mp3&LogLevel=ERROR"
	if code := alert.StatusCode(); code != 405 {
		t.Errorf("expected Code to be 405, got %d", code)
	}
	alert.AlertText = "Msg&sourceComponent=12000&ErrorCode=11200&httpResponse=4050&url=https%3A%2F%2Fkev.inburke.com%2Fzombo%2Fzombocom.mp3&LogLevel=ERROR"
	if code := alert.StatusCode(); code != 0 {
		t.Errorf("expected Code to be 0, got %d", code)
	}
}
