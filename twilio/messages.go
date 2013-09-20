package twilio

import (
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

func (m *MessageService) Create(data url.Values) (Message, error) {
	msg := new(Message)
	err := m.client.MakeRequest("POST", pathPart, data, msg)
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
