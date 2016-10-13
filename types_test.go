package twilio

import (
	"encoding/json"
	"testing"
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
	for _, tt := range pnTestCases {
		if f := tt.in.Friendly(); f != tt.expected {
			t.Errorf("Friendly(%v): got %s, want %s", tt.in, f, tt.expected)
		}
	}
}

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
	for _, tt := range priceTests {
		out := price(tt.unit, tt.amount)
		if out != tt.expected {
			t.Errorf("price(%v, %v): got %v, want %v", tt.unit, tt.amount, out, tt.expected)
		}
	}
}
