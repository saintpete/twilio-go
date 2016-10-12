package twilio

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalTime(t *testing.T) {
	in := []byte(`"Tue, 20 Sep 2016 22:59:57 +0000"`)
	var tt TwilioTime
	if err := json.Unmarshal(in, &tt); err != nil {
		t.Fatal(err)
	}
	if tt.Valid == false {
		t.Errorf("expected time to be Valid, got false")
	}
	if tt.Time.Year() != 2016 {
		t.Errorf("expected Year to equal 2016, got %d", tt.Time.Year())
	}
}
