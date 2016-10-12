package twilio

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ttacon/libphonenumber"
)

type PhoneNumber string

// Friendly returns a friendly international representation of the phone
// number, for example, "+14105554092" is returned as "+1 410-555-4092". If the
// phone number is not in E.164 format, we try to parse it as a US number. If
// we cannot parse it as a US number, it is returned as is.
func (pn PhoneNumber) Friendly() string {
	num, err := libphonenumber.Parse(string(pn), "US")
	if err != nil {
		return string(pn)
	}
	return libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
}

// Local returns a friendly national representation of the phone number, for
// example, "+14105554092" is returned as "(410) 555-4092". If the phone number
// is not in E.164 format, we try to parse it as a US number. If we cannot
// parse it as a US number, it is returned as is.
func (pn PhoneNumber) Local() string {
	num, err := libphonenumber.Parse(string(pn), "US")
	if err != nil {
		return string(pn)
	}
	return libphonenumber.Format(num, libphonenumber.NATIONAL)
}

type Segments uint

func (seg *Segments) UnmarshalJSON(b []byte) error {
	s := new(string)
	if err := json.Unmarshal(b, s); err != nil {
		return err
	}
	u, err := strconv.ParseUint(*s, 10, 64)
	if err != nil {
		return err
	}
	*seg = Segments(u)
	return nil
}

// TwilioTime can parse a timestamp returned in the Twilio API and turn it into
// a valid Go Time struct.
type TwilioTime struct {
	Time  time.Time
	Valid bool
}

// The reference time, as it appears in the Twilio API.
const TimeLayout = "Mon, 2 Jan 2006 15:04:05 -0700"

func (t *TwilioTime) UnmarshalJSON(b []byte) error {
	s := new(string)
	if err := json.Unmarshal(b, s); err != nil {
		return err
	}
	if *s == "null" {
		t.Valid = false
		return nil
	}
	tim, err := time.Parse(TimeLayout, *s)
	if err != nil {
		return err
	}
	*t = TwilioTime{Time: tim, Valid: true}
	return nil
}

func (tt *TwilioTime) MarshalJSON() ([]byte, error) {
	if tt.Valid == false {
		return []byte("null"), nil
	}
	b, err := json.Marshal(tt.Time)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
