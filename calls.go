package twilio

import (
	"encoding/json"
	"net/url"

	types "github.com/kevinburke/go-types"
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

func (c *CallService) Get(sid string) (*Call, error) {
	call := new(Call)
	err := c.client.GetResource(callsPathPart, sid, call)
	return call, err
}

func (c *CallService) GetPage(data url.Values) (*CallPage, error) {
	cp := new(CallPage)
	err := c.client.ListResource(callsPathPart, data, cp)
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
func (c *CallPageIterator) Next() (*CallPage, error) {
	cp := new(CallPage)
	err := c.p.Next(cp)
	if err != nil {
		return nil, err
	}
	c.p.SetNextPageURI(cp.NextPageURI)
	return cp, nil
}
