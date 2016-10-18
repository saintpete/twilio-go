package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"
)

var pnTestCases = []struct {
	in       PhoneNumber
	expected string
}{
	{PhoneNumber("+41446681800"), "+41 44 668 18 00"},
	{PhoneNumber("+14105554092"), "+1 410-555-4092"},
	{PhoneNumber("blah"), "blah"},
}

func TestPhoneNumberFriendly(t *testing.T) {
	t.Parallel()
	for _, tt := range pnTestCases {
		if f := tt.in.Friendly(); f != tt.expected {
			t.Errorf("Friendly(%v): got %s, want %s", tt.in, f, tt.expected)
		}
	}
}

var pnParseTestCases = []struct {
	in  string
	out PhoneNumber
	err error
}{
	{"+14105551234", PhoneNumber("+14105551234"), nil},
	{"410 555 1234", PhoneNumber("+14105551234"), nil},
	{"(410) 555-1234", PhoneNumber("+14105551234"), nil},
	{"+41 44 6681800", PhoneNumber("+41446681800"), nil},
	{"foobarbang", PhoneNumber(""), errors.New("twilio: Invalid phone number: foobarbang")},
	{"22", PhoneNumber("+122"), nil},
	{"", PhoneNumber(""), ErrEmptyNumber},
}

func TestNewPhoneNumber(t *testing.T) {
	t.Parallel()
	for _, tt := range pnParseTestCases {
		pn, err := NewPhoneNumber(tt.in)
		name := fmt.Sprintf("ParsePhoneNumber(%v)", tt.in)
		if tt.err != nil {
			if err == nil {
				t.Errorf("%s: expected %v, got nil", name, tt.err)
				continue
			}
			if err.Error() != tt.err.Error() {
				t.Errorf("%s: expected error %v, got %v", name, tt.err, err)
			}
		} else if pn != tt.out {
			t.Errorf("%s: expected %v, got %v", name, tt.out, pn)
		}
	}
}

func TestUnmarshalTime(t *testing.T) {
	t.Parallel()
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
	in = []byte(`null`)
	if err := json.Unmarshal(in, &tt); err != nil {
		t.Fatal(err)
	}
	if tt.Valid != false {
		t.Errorf("expected time.Valid to be false, got true")
	}
}

func TestNewTwilioTime(t *testing.T) {
	t.Parallel()
	v := NewTwilioTime("foo")
	if v.Valid == true {
		t.Errorf("expected time to be invalid, got true")
	}
	in := "Tue, 20 Sep 2016 22:59:57 +0000"
	v = NewTwilioTime(in)
	if v.Valid == false {
		t.Errorf("expected %s to be valid time, got false", in)
	}
	if v.Time.Minute() != 59 {
		t.Errorf("wrong minute")
	}
	expected := "2016-09-20T22:59:57Z"
	if f := v.Time.Format(time.RFC3339); f != expected {
		t.Errorf("wrong format: got %v, want %v", f, expected)
	}
}

var priceTests = []struct {
	unit     string
	amount   string
	expected string
}{
	{"USD", "-0.0075", "$0.0075"},
	{"usd", "-0.0075", "$0.0075"},
	{"EUR", "0.0075", "€-0.0075"},
	{"UNK", "2.45", "UNK -2.45"},
	{"", "2.45", "-2.45"},
	{"USD", "-0.75000000", "$0.75"},
	{"USD", "-0.750", "$0.75"},
	{"USD", "-5000.00", "$5000"},
	{"USD", "-5000.", "$5000"},
	{"USD", "-5000", "$5000"},
}

func TestPrice(t *testing.T) {
	t.Parallel()
	for _, tt := range priceTests {
		out := price(tt.unit, tt.amount)
		if out != tt.expected {
			t.Errorf("price(%v, %v): got %v, want %v", tt.unit, tt.amount, out, tt.expected)
		}
	}
}

func TestTwilioDuration(t *testing.T) {
	t.Parallel()
	in := []byte(`"88"`)
	var td TwilioDuration
	if err := json.Unmarshal(in, &td); err != nil {
		t.Fatal(err)
	}
	if td != TwilioDuration(88*time.Second) {
		t.Errorf("got wrong duration: %v, wanted 88 seconds", td)
	}
}
