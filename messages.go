package twilio

import (
	"net/url"

	types "github.com/Shyp/go-types"
)

const pathPart = "Messages"

type MessageService struct {
	client *Client
}

type Direction string

const DirectionOutboundReply = Direction("outbound-reply")

type Status string

const StatusSent = Status("sent")
const StatusReceived = Status("received")
const StatusDelivered = Status("delivered")

type Message struct {
	Sid                 string           `json:"sid"`
	Body                string           `json:"body"`
	From                string           `json:"from"`
	To                  string           `json:"to"`
	Price               string           `json:"price"`
	Status              Status           `json:"status"`
	AccountSid          string           `json:"account_sid"`
	MessagingServiceSid types.NullString `json:"messaging_service_sid"`
	DateCreated         TwilioTime       `json:"date_created"`
	DateUpdated         TwilioTime       `json:"date_updated"`
	DateSent            TwilioTime       `json:"date_sent"`
	NumSegments         Segments         `json:"num_segments"`
	// TODO fix type here... UintStr or something ?
	NumMedia        Segments          `json:"num_media"`
	PriceUnit       string            `json:"price_unit"`
	Direction       Direction         `json:"direction"`
	SubresourceURIs map[string]string `json:"subresource_uris"`
	URI             string            `json:"uri"`
	APIVersion      string            `json:"api_version"`
}

type MessagePage struct {
	FirstPageURI string     `json:"first_page_uri"`
	Start        uint       `json:"start"`
	End          uint       `json:"end"`
	NumPages     uint       `json:"num_pages"`
	Total        uint       `json:"total"`
	NextPageURI  string     `json:"next_page_uri"`
	PageSize     uint       `json:"page_size"`
	Messages     []*Message `json:"messages"`
}

type MessageIterator struct {
	pos         int
	messages    []*Message
	nextPageURI string
	data        url.Values
	client      *Client
}

// Create a message with the given values.
func (m *MessageService) Create(data url.Values) (*Message, error) {
	msg := new(Message)
	err := m.client.CreateResource(pathPart, data, msg)
	return msg, err
}

// SendMessage is a convenience wrapper around Create.
func (m *MessageService) SendMessage(from string, to string, body string, mediaURLs []url.URL) (*Message, error) {
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

// List returns a MessageIteratar with the given values.
func (m *MessageService) List(data url.Values) *MessageIterator {
	return &MessageIterator{
		pos:    0,
		data:   data,
		client: m.client,
	}
}

// Next gets the next message. Returns nil, io.EOF when there are no more
// messages to read.
func (m *MessageIterator) Next() (*Message, error) {
	if m.pos >= len(m.messages) {
		if m.nextPageURI != "" {
			m.data.Set("next_page_uri", m.nextPageURI)
		}
		page := new(MessagePage)
		err := m.client.ListResource(pathPart, m.data, &page)
		if err != nil {
			return nil, err
		}
		m.nextPageURI = page.NextPageURI
		m.messages = page.Messages
		m.pos = 0
	}
	m.pos += 1
	return m.messages[m.pos-1], nil
}

// GetPage returns a single page of messages.
func (m *MessageService) GetPage(data url.Values) (*MessagePage, error) {
	page := new(MessagePage)
	err := m.client.ListResource(pathPart, data, page)
	return page, err
}
