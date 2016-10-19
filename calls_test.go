package twilio

import "testing"

func TestGetCall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "CAa98f7bbc9bc4980a44b128ca4884ca73"
	call, err := envClient.Calls.Get(sid)
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
	recordings, err := envClient.Calls.GetRecordings(sid, nil)
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
