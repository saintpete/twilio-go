package twilio

import (
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGetCall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "CAa98f7bbc9bc4980a44b128ca4884ca73"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	call, err := envClient.Calls.Get(ctx, sid)
	if err != nil {
		t.Fatal(err)
	}
	if call.Sid != sid {
		t.Errorf("expected Sid to equal %s, got %s", sid, call.Sid)
	}
}

func TestGetCallRecordings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "CA14365760c10f73392c5440bdfb70c212"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	recordings, err := envClient.Calls.GetRecordings(ctx, sid, nil)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(recordings.Recordings); l != 1 {
		t.Fatalf("expected 1 recording, got %d", l)
	}
	rsid := "REd04242a0544234abba080942e0535505"
	if r := recordings.Recordings[0].Sid; r != rsid {
		t.Errorf("expected recording sid to be %s, got %s", rsid, r)
	}
	if recordings.NextPageURI.Valid {
		t.Errorf("expected next page uri to be invalid, got %v", recordings.NextPageURI)
	}
}

func TestMakeCall(t *testing.T) {
	t.Parallel()
	client, server := getServer(makeCallResponse)
	defer server.Close()
	u, _ := url.Parse("https://kev.inburke.com/zombo/zombocom.mp3")
	call, err := client.Calls.MakeCall("+19253920364", "+19252717005", u)
	if err != nil {
		t.Fatal(err)
	}
	if call.To != "+19252717005" {
		t.Errorf("Wrong To phone number: %s", call.To)
	}
	if call.Status != StatusQueued {
		t.Errorf("Wrong status: %s", call.Status)
	}
}

func TestGetCallRange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	// from Kevin's account, there should be exactly 2 results in this range.
	data := url.Values{}
	data.Set("PageSize", "2")
	nyc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	// 10:34:00 Oct 26 to 19:25:59 Oct 27 NYC time. I made 2 calls in this
	// range, and 5 calls on this day. There are 2 calls before this time on
	// this day, and 1 call after.
	start := time.Date(2016, 10, 26, 22, 34, 00, 00, nyc)
	end := time.Date(2016, 10, 27, 19, 25, 59, 00, nyc)
	iter := envClient.Calls.GetCallsInRange(start, end, data)
	count := 0
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		count++
		page, err := iter.Next(ctx)
		if err == NoMoreResults {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Calls) != 2 {
			t.Errorf("expected 2 calls in result set, got %d", len(page.Calls))
			break
		}
		// the 19:25 call on 10-27
		if page.Calls[0].Sid != "CA5757109d6dcbc4bebf6847f5dd45191e" {
			t.Errorf("wrong sid")
		}
		// the 22:34 call on 10-26
		if page.Calls[1].Sid != "CA47b862ce3b99a6d79939320a9aa54a02" {
			t.Errorf("wrong sid")
		}
	}
	if count != 2 {
		t.Errorf("wrong count, expected exactly 2 calls to Next(), got %d", count)
	}
}
