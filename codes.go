package twilio

import (
	"encoding/json"
	"strconv"
)

// A Twilio error code. A full list can be found here:
// https://www.twilio.com/docs/api/errors/reference
type Code int

func (c *Code) UnmarshalJSON(b []byte) error {
	s := new(string)
	if err := json.Unmarshal(b, s); err == nil {
		if *s == "" || *s == "null" {
			*c = Code(0)
			return nil
		}
		i, err := strconv.Atoi(*s)
		if err != nil {
			return err
		}
		*c = Code(i)
		return nil
	}
	i := new(int)
	err := json.Unmarshal(b, i)
	if err != nil {
		return err
	}
	*c = Code(*i)
	return nil
}

const CodeHTTPRetrievalFailure = 11200
const CodeHTTPConnectionFailure = 11205
const CodeHTTPProtocolViolation = 11206
const CodeQueueOverflow = 30001
const CodeAccountSuspended = 30002
const CodeUnreachable = 30003
const CodeMessageBlocked = 30004
const CodeUnknownDestination = 30005
const CodeLandline = 30006
const CodeCarrierViolation = 30007
const CodeUnknownError = 30008
const CodeMissingSegment = 30009
const CodeMessagePriceExceedsMaxPrice = 30010
