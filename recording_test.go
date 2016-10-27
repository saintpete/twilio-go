package twilio

import (
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGetRecording(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "REd04242a0544234abba080942e0535505"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	recording, err := envClient.Recordings.Get(ctx, sid)
	if err != nil {
		t.Fatal(err)
	}
	if recording.Sid != sid {
		t.Errorf("expected Sid to equal %s, got %s", sid, recording.Sid)
	}
}
