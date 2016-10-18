package twilio

// A Twilio error code. A full list can be found here:
// https://www.twilio.com/docs/api/errors/reference
type Code int

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
