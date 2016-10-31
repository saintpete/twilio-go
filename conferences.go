package twilio

import (
	"net/url"

	"golang.org/x/net/context"
)

const conferencePathPart = "Conferences"

type ConferenceService struct {
	client *Client
}

type Conference struct {
	Sid string `json:"sid"`
	// Call status, StatusInProgress or StatusCompleted
	Status       Status `json:"status"`
	FriendlyName string `json:"friendly_name"`
	// The conference region, probably "us1"
	Region      string     `json:"region"`
	DateCreated TwilioTime `json:"date_created"`
	AccountSid  string     `json:"account_sid"`
	APIVersion  string     `json:"api_version"`
	DateUpdated TwilioTime `json:"date_updated"`
	URI         string     `json:"uri"`
}

type ConferencePage struct {
	Page
	Conferences []*Conference
}

func (c *ConferenceService) Get(ctx context.Context, sid string) (*Conference, error) {
	conference := new(Conference)
	err := c.client.GetResource(ctx, conferencePathPart, sid, conference)
	return conference, err
}

func (c *ConferenceService) GetPage(ctx context.Context, data url.Values) (*ConferencePage, error) {
	cp := new(ConferencePage)
	err := c.client.ListResource(ctx, conferencePathPart, data, cp)
	return cp, err
}

// ConferencePageIterator lets you retrieve consecutive ConferencePages.
type ConferencePageIterator struct {
	p *PageIterator
}

// GetPageIterator returns a ConferencePageIterator with the given page
// filters. Call iterator.Next() to get the first page of resources (and again
// to retrieve subsequent pages).
func (c *ConferenceService) GetPageIterator(data url.Values) *ConferencePageIterator {
	iter := NewPageIterator(c.client, data, conferencePathPart)
	return &ConferencePageIterator{
		p: iter,
	}
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (c *ConferencePageIterator) Next(ctx context.Context) (*ConferencePage, error) {
	cp := new(ConferencePage)
	err := c.p.Next(ctx, cp)
	if err != nil {
		return nil, err
	}
	c.p.SetNextPageURI(cp.NextPageURI)
	return cp, nil
}
