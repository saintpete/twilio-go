package twilio

import "testing"

func TestGetRecording(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "REd04242a0544234abba080942e0535505"
	recording, err := envClient.Recordings.Get(sid)
	if err != nil {
		t.Fatal(err)
	}
	if recording.Sid != sid {
		t.Errorf("expected Sid to equal %s, got %s", sid, recording.Sid)
	}
}
