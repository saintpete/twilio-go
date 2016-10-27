package twilio

import (
	"encoding/json"
	"net/url"

	types "github.com/kevinburke/go-types"
	"golang.org/x/net/context"
)

const callsPathPart = "Calls"

type CallService struct {
	client *Client
}

type Call struct {
	Sid            string           `json:"sid"`
	From           PhoneNumber      `json:"from"`
	To             PhoneNumber      `json:"to"`
	Status         Status           `json:"status"`
	StartTime      TwilioTime       `json:"start_time"`
	EndTime        TwilioTime       `json:"end_time"`
	Duration       TwilioDuration   `json:"duration"`
	AccountSid     string           `json:"account_sid"`
	Annotation     json.RawMessage  `json:"annotation"`
	AnsweredBy     NullAnsweredBy   `json:"answered_by"`
	CallerName     types.NullString `json:"caller_name"`
	DateCreated    TwilioTime       `json:"date_created"`
	DateUpdated    TwilioTime       `json:"date_updated"`
	Direction      Direction        `json:"direction"`
	ForwardedFrom  PhoneNumber      `json:"forwarded_from"`
	GroupSid       string           `json:"group_sid"`
	ParentCallSid  string           `json:"parent_call_sid"`
	PhoneNumberSid string           `json:"phone_number_sid"`
	Price          string           `json:"price"`
	PriceUnit      string           `json:"price_unit"`
	APIVersion     string           `json:"api_version"`
	URI            string           `json:"uri"`
}

// FriendlyPrice flips the sign of the Price (which is usually reported from
// the API as a negative number) and adds an appropriate currency symbol in
// front of it. For example, a PriceUnit of "USD" and a Price of "-1.25" is
// reported as "$1.25".
func (c *Call) FriendlyPrice() string {
	if c == nil {
		return ""
	}
	return price(c.PriceUnit, c.Price)
}

// A CallPage contains a Page of calls.
type CallPage struct {
	Page
	Calls []*Call `json:"calls"`
}

func (c *CallService) Get(ctx context.Context, sid string) (*Call, error) {
	call := new(Call)
	err := c.client.GetResource(ctx, callsPathPart, sid, call)
	return call, err
}

// Update the call with the given data. Valid parameters may be found here:
// https://www.twilio.com/docs/api/rest/change-call-state#post-parameters
func (c *CallService) Update(ctx context.Context, sid string, data url.Values) (*Call, error) {
	call := new(Call)
	err := c.client.UpdateResource(ctx, callsPathPart, sid, data, call)
	return call, err
}

// Cancel an in-progress Call with the given sid. Cancel will not affect
// in-progress Calls, only those in queued or ringing.
func (c *CallService) Cancel(sid string) (*Call, error) {
	data := url.Values{}
	data.Set("Status", string(StatusCanceled))
	return c.Update(context.Background(), sid, data)
}

// Hang up an in-progress call.
func (c *CallService) Hangup(sid string) (*Call, error) {
	data := url.Values{}
	data.Set("Status", string(StatusCompleted))
	return c.Update(context.Background(), sid, data)
}

// Redirect the given call to the given URL.
func (c *CallService) Redirect(sid string, u *url.URL) (*Call, error) {
	data := url.Values{}
	data.Set("Url", u.String())
	return c.Update(context.Background(), sid, data)
}

// Initiate a new Call.
func (c *CallService) Create(ctx context.Context, data url.Values) (*Call, error) {
	call := new(Call)
	err := c.client.CreateResource(ctx, callsPathPart, data, call)
	return call, err
}

// MakeCall starts a new Call from the given phone number to the given phone
// number, dialing the url when the call connects. MakeCall is a wrapper around
// Create; if you need more configuration, call that function directly.
func (c *CallService) MakeCall(from string, to string, u *url.URL) (*Call, error) {
	data := url.Values{}
	data.Set("From", from)
	data.Set("To", to)
	data.Set("Url", u.String())
	return c.Create(context.Background(), data)
}

func (c *CallService) GetPage(ctx context.Context, data url.Values) (*CallPage, error) {
	cp := new(CallPage)
	err := c.client.ListResource(ctx, callsPathPart, data, cp)
	return cp, err
}

type CallPageIterator struct {
	p *PageIterator
}

// GetPageIterator returns an iterator which can be used to retrieve pages.
func (c *CallService) GetPageIterator(data url.Values) *CallPageIterator {
	iter := NewPageIterator(c.client, data, callsPathPart)
	return &CallPageIterator{
		p: iter,
	}
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (c *CallPageIterator) Next(ctx context.Context) (*CallPage, error) {
	cp := new(CallPage)
	err := c.p.Next(ctx, cp)
	if err != nil {
		return nil, err
	}
	c.p.SetNextPageURI(cp.NextPageURI)
	return cp, nil
}

// GetRecordings returns an array of recordings for this Call. Note there may
// be more than one Page of results.
func (c *CallService) GetRecordings(ctx context.Context, callSid string, data url.Values) (*RecordingPage, error) {
	if data == nil {
		data = url.Values{}
	}
	// Cheat - hit the Recordings list view with a filter instead of
	// GET /calls/CA123/Recordings. The former is probably more reliable
	data.Set("CallSid", callSid)
	return c.client.Recordings.GetPage(ctx, data)
}

// GetRecordings returns an iterator of recording pages for this Call.
// Note there may be more than one Page of results.
func (c *CallService) GetRecordingsIterator(callSid string, data url.Values) *RecordingPageIterator {
	if data == nil {
		data = url.Values{}
	}
	// Cheat - hit the Recordings list view with a filter instead of
	// GET /calls/CA123/Recordings. The former is probably more reliable
	data.Set("CallSid", callSid)
	return c.client.Recordings.GetPageIterator(data)
}
