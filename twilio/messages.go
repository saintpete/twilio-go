package twilio

import (
	//"fmt"
	"github.com/golang/glog"
	"net/url"
)

const pathPart = "Messages"

type MessageService struct {
	client *Client
}

type Message struct {
	Body  string
	From  string
	To    string
	Price string
	Sid   string
}

type MessagePage struct {
	NextPageUri string    `json:"next_page_uri"`
	PageSize    int       `json:"page_size"`
	Messages    []Message `json:"messages"`
}

type MessageIterator struct {
	pos         int
	messages    []Message
	nextPageUri string
	data        url.Values
	client      *Client
}

func (m *MessageService) Create(data url.Values) (Message, error) {
	msg := new(Message)
	_, err := m.client.MakeRequest("POST", pathPart, data, msg)
	if err != nil {
		glog.Errorf("Error creating request", err)
		return *msg, err
	}
	return *msg, nil
}

func (m *MessageService) SendMessage(from string, to string, body string, mediaUrls []url.URL) (Message, error) {
	v := url.Values{}
	v.Set("Body", body)
	v.Set("From", from)
	v.Set("To", to)
	if mediaUrls != nil {
		for _, mediaUrl := range mediaUrls {
			v.Add("MediaUrl", mediaUrl.String())
		}
	}
	return m.Create(v)
}

func (m *MessageService) ListMessages(data url.Values) MessageIterator {
	iterator := new(MessageIterator)
	if data == nil {
		iterator.data = url.Values{
			"Page": {"5000"},
		}
	} else {
		iterator.data = data
	}
	iterator.pos = 2000
	iterator.client = m.client
	return *iterator
}

// Returns nil when the list is complete.
func (m *MessageIterator) Next() (*Message, error) {
	if m.pos >= len(m.messages) {
		if m.nextPageUri != "" {
			m.data.Set("next_page_uri", m.nextPageUri)
		}
		var page MessagePage
		_, err := m.client.ListResource(pathPart, m.data, &page)
		if err != nil {
			glog.Errorf("Error creating request", err)
			return nil, err
		}
		m.nextPageUri = page.NextPageUri
		m.messages = page.Messages
		m.pos = 0
	}
	m.pos += 1
	return &m.messages[m.pos-1], nil
}
