package twilio

import (
	"net/url"
	"time"

	types "github.com/Shyp/go-types"
)

const pathPart = "Messages"

type MessageService struct {
	client *Client
}

type Message struct {
	Sid         string         `json:"sid"`
	Body        string         `json:"body"`
	From        string         `json:"from"`
	To          string         `json:"to"`
	Price       string         `json:"price"`
	DateCreated time.Time      `json:"date_created"`
	DateUpdated time.Time      `json:"date_updated"`
	DateSent    types.NullTime `json:"date_sent"`
	NumSegments uint           `json:"num_segments"`
}

type MessagePage struct {
	NextPageUri string     `json:"next_page_uri"`
	PageSize    uint       `json:"page_size"`
	Messages    []*Message `json:"messages"`
}

type MessageIterator struct {
	pos         int
	messages    []*Message
	nextPageUri string
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
		if m.nextPageUri != "" {
			m.data.Set("next_page_uri", m.nextPageUri)
		}
		page := new(MessagePage)
		err := m.client.ListResource(pathPart, m.data, &page)
		if err != nil {
			return nil, err
		}
		m.nextPageUri = page.NextPageUri
		m.messages = page.Messages
		m.pos = 0
	}
	m.pos += 1
	return m.messages[m.pos-1], nil
}
