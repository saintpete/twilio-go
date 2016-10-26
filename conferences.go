package twilio

import "net/url"

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

func (c *ConferenceService) Get(sid string) (*Conference, error) {
	conference := new(Conference)
	err := c.client.GetResource(conferencePathPart, sid, conference)
	return conference, err
}

func (c *ConferenceService) GetPage(data url.Values) (*ConferencePage, error) {
	cp := new(ConferencePage)
	err := c.client.ListResource(conferencePathPart, data, cp)
	return cp, err
}

type ConferencePageIterator struct {
	p *PageIterator
}

func (c *ConferenceService) GetPageIterator(data url.Values) *ConferencePageIterator {
	iter := NewPageIterator(c.client, data, conferencePathPart)
	return &ConferencePageIterator{
		p: iter,
	}
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (c *ConferencePageIterator) Next() (*ConferencePage, error) {
	cp := new(ConferencePage)
	err := c.p.Next(cp)
	if err != nil {
		return nil, err
	}
	c.p.SetNextPageURI(cp.NextPageURI)
	return cp, nil
}
