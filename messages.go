package twilio

import (
	"net/url"
	"strings"
	"sync"

	types "github.com/kevinburke/go-types"
)

const pathPart = "Messages"

type MessageService struct {
	client *Client
}

// The direction of the message.
type Direction string

// Friendly prints out a friendly version of the Direction, following the
// example shown in the Twilio Dashboard.
func (d Direction) Friendly() string {
	switch {
	case d == DirectionOutboundReply:
		return "Reply"
	case d == DirectionOutboundCall:
		return "Outgoing (from call)"
	case d == DirectionOutboundAPI:
		return "Outgoing (from API)"
	case d == DirectionInbound:
		return "Incoming"
	default:
		return string(d)
	}
}

const DirectionOutboundReply = Direction("outbound-reply")
const DirectionInbound = Direction("inbound")
const DirectionOutboundCall = Direction("outbound-call")
const DirectionOutboundAPI = Direction("outbound-api")

// The status of the message (accepted, queued, etc).
// For more information , see https://www.twilio.com/docs/api/rest/message
type Status string

func (s Status) Friendly() string {
	return strings.Title(string(s))
}

const StatusAccepted = Status("accepted")
const StatusDelivered = Status("delivered")
const StatusFailed = Status("failed")
const StatusQueued = Status("queued")
const StatusReceiving = Status("receiving")
const StatusReceived = Status("received")
const StatusSending = Status("sending")
const StatusSent = Status("sent")
const StatusUndelivered = Status("undelivered")

type Message struct {
	Sid                 string            `json:"sid"`
	Body                string            `json:"body"`
	From                PhoneNumber       `json:"from"`
	To                  PhoneNumber       `json:"to"`
	Price               string            `json:"price"`
	Status              Status            `json:"status"`
	AccountSid          string            `json:"account_sid"`
	MessagingServiceSid types.NullString  `json:"messaging_service_sid"`
	DateCreated         TwilioTime        `json:"date_created"`
	DateUpdated         TwilioTime        `json:"date_updated"`
	DateSent            TwilioTime        `json:"date_sent"`
	NumSegments         Segments          `json:"num_segments"`
	NumMedia            NumMedia          `json:"num_media"`
	PriceUnit           string            `json:"price_unit"`
	Direction           Direction         `json:"direction"`
	SubresourceURIs     map[string]string `json:"subresource_uris"`
	URI                 string            `json:"uri"`
	APIVersion          string            `json:"api_version"`
}

// FriendlyPrice flips the sign of the Price (which is usually reported from
// the API as a negative number) and adds an appropriate currency symbol in
// front of it. For example, a PriceUnit of "USD" and a Price of "-1.25" is
// reported as "$1.25".
func (m *Message) FriendlyPrice() string {
	return price(m.PriceUnit, m.Price)
}

// A MessagePage contains a Page of messages.
type MessagePage struct {
	Page
	Messages []*Message `json:"messages"`
}

// Create a message with the given url.Values. For more information on valid
// values, see https://www.twilio.com/docs/api/rest/sending-messages or use the
// SendMessage helper.
func (m *MessageService) Create(data url.Values) (*Message, error) {
	msg := new(Message)
	err := m.client.CreateResource(pathPart, data, msg)
	return msg, err
}

// SendMessage is a convenience wrapper around Create.
func (m *MessageService) SendMessage(from string, to string, body string, mediaURLs []*url.URL) (*Message, error) {
	v := url.Values{
		"Body": []string{body},
		"From": []string{from},
		"To":   []string{to},
	}
	if mediaURLs != nil {
		for _, mediaURL := range mediaURLs {
			v.Add("MediaUrl", mediaURL.String())
		}
	}
	return m.Create(v)
}

type MessagePageIterator struct {
	client      *Client
	nextPageURI types.NullString
	data        url.Values
	count       uint
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (m *MessagePageIterator) Next() (*MessagePage, error) {
	mp := new(MessagePage)
	var err error
	if m.count == 0 {
		err = m.client.ListResource(pathPart, m.data, mp)
	} else if m.nextPageURI.Valid == false {
		return nil, NoMoreResults
	} else {
		err = m.client.GetNextPage(m.nextPageURI.String, mp)
	}
	if err != nil {
		return nil, err
	}
	m.count++
	m.nextPageURI = mp.NextPageURI
	return mp, nil
}

func (m *MessageService) Get(sid string) (*Message, error) {
	msg := new(Message)
	err := m.client.GetResource(pathPart, sid, msg)
	return msg, err
}

// GetPage returns a single page of resources. To retrieve multiple pages, use
// GetPageIterator.
func (m *MessageService) GetPage(data url.Values) (*MessagePage, error) {
	mp := new(MessagePage)
	err := m.client.ListResource(pathPart, data, mp)
	return mp, err
}

// GetPageIterator returns an iterator which can be used to retrieve pages.
func (m *MessageService) GetPageIterator(data url.Values) *MessagePageIterator {
	return &MessagePageIterator{
		client:      m.client,
		nextPageURI: types.NullString{},
		data:        data,
		count:       0,
	}
}

// GetMediaURLs gets the URLs of any media for this message. This uses threads
// to retrieve all URLs simultaneously; if retrieving any URL fails, we return
// an error for the entire request.
//
// The data can be used to filter the list of returned Media as described here:
// https://www.twilio.com/docs/api/rest/media#list-get-filters
//
// As of October 2016, only 10 MediaURLs are permitted per message. No attempt
// is made to page through media resources; omit the PageSize parameter in
// data, or set it to a value greater than 10, to retrieve all resources.
func (m *MessageService) GetMediaURLs(sid string, data url.Values) ([]*url.URL, error) {
	page, err := m.client.Media.GetPage(sid, data)
	if err != nil {
		return nil, err
	}
	if len(page.MediaList) == 0 {
		urls := make([]*url.URL, 0, 0)
		return urls, nil
	}
	urls := make([]*url.URL, len(page.MediaList))
	errs := make([]error, len(page.MediaList))
	var wg sync.WaitGroup
	wg.Add(len(page.MediaList))
	for i, media := range page.MediaList {
		go func(i int, media *Media) {
			url, err := m.client.Media.GetURL(sid, media.Sid)
			urls[i] = url
			errs[i] = err
			wg.Done()
		}(i, media)
	}
	wg.Wait()
	// todo - we could probably return more quickly in the result of a failure.
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return urls, nil
}
