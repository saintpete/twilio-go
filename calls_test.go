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
