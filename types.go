package twilio

import (
	"encoding/json"
	"strconv"
	"time"
)

type Segments uint

func (t *Segments) UnmarshalJSON(b []byte) error {
	s := new(string)
	if err := json.Unmarshal(b, s); err != nil {
		return err
	}
	u, err := strconv.ParseUint(*s, 10, 64)
	if err != nil {
		return err
	}
	*t = Segments(u)
	return nil
}

type TwilioTime struct {
	Time  time.Time
	Valid bool
}

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
