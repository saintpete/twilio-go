package twilio

import (
	"context"
	"net/url"

	types "github.com/kevinburke/go-types"
)

const simPathPart = "Sims"

type SimService struct {
	client *Client
}

// Sim represents a Sim resource.
type Sim struct {
	Sid          string           `json:"sid"`
	UniqueName   string           `json:"unique_name"`
	Status       Status           `json:"status"`
	FriendlyName types.NullString `json:"friendly_name"`
	ICCID        string           `json:"iccid"`

	CommandsCallbackMethod string           `json:"commands_callback_method"`
	CommandsCallbackURL    types.NullString `json:"commands_callback_url"`
	DateCreated            TwilioTime       `json:"date_created"`
	DateUpdated            TwilioTime       `json:"date_updated"`
	RatePlanSid            string           `json:"rate_plan_sid"`
	SMSURL                 types.NullString `json:"sms_url"`
	SMSMethod              types.NullString `json:"sms_method"`
	SMSFallbackMethod      types.NullString `json:"sms_fallback_method"`
	SMSFallbackURL         types.NullString `json:"sms_fallback_url"`
	VoiceURL               types.NullString `json:"voice_url"`
	VoiceMethod            types.NullString `json:"voice_method"`
	VoiceFallbackMethod    types.NullString `json:"voice_fallback_method"`
	VoiceFallbackURL       types.NullString `json:"voice_fallback_url"`

	URL        string            `json:"url"`
	AccountSid string            `json:"account_sid"`
	Links      map[string]string `json:"links"`
}

// SimPage represents a page of Sims.
type SimPage struct {
	Meta Meta   `json:"meta"`
	Sims []*Sim `json:"sims"`
}

// Get finds a single Sim resource by its sid, or returns an error.
func (s *SimService) Get(ctx context.Context, sid string) (*Sim, error) {
	sim := new(Sim)
	err := s.client.GetResource(ctx, simPathPart, sid, sim)
	return sim, err
}

// Update the sim with the given data. Valid parameters may be found here:
// https://www.twilio.com/docs/api/wireless/rest-api/sim#instance-post
func (c *SimService) Update(ctx context.Context, sid string, data url.Values) (*Sim, error) {
	sim := new(Sim)
	err := c.client.UpdateResource(ctx, simPathPart, sid, data, sim)
	return sim, err
}

// SimPageIterator lets you retrieve consecutive pages of resources.
type SimPageIterator interface {
	// Next returns the next page of resources. If there are no more resources,
	// NoMoreResults is returned.
	Next(context.Context) (*SimPage, error)
}

type simPageIterator struct {
	p *PageIterator
}

// GetPage returns a single Page of resources, filtered by data.
//
// See https://www.twilio.com/docs/api/wireless/rest-api/sim#list-get.
func (f *SimService) GetPage(ctx context.Context, data url.Values) (*SimPage, error) {
	return f.GetPageIterator(data).Next(ctx)
}

// GetPageIterator returns a SimPageIterator with the given page
// filters. Call iterator.Next() to get the first page of resources (and again
// to retrieve subsequent pages).
func (f *SimService) GetPageIterator(data url.Values) SimPageIterator {
	iter := NewPageIterator(f.client, data, simPathPart)
	return &simPageIterator{
		p: iter,
	}
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (f *simPageIterator) Next(ctx context.Context) (*SimPage, error) {
	ap := new(SimPage)
	err := f.p.Next(ctx, ap)
	if err != nil {
		return nil, err
	}
	f.p.SetNextPageURI(ap.Meta.NextPageURL)
	return ap, nil
}
