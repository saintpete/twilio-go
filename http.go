package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kevinburke/rest"
)

const BaseURL = "https://api.twilio.com"
const APIVersion = "2010-04-01"

type Client struct {
	*rest.Client

	AccountSid string
	AuthToken  string

	Messages *MessageService
}

const defaultTimeout = 30*time.Second + 500*time.Millisecond

// NewClient creates a Client for interacting with the Twilio API.
func NewClient(accountSid string, authToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	restClient := rest.NewClient(accountSid, authToken, fmt.Sprintf("%s/%s", BaseURL, APIVersion))
	restClient.Client = httpClient

	c := &Client{Client: restClient, AccountSid: accountSid, AuthToken: authToken}
	c.Messages = &MessageService{client: c}
	return c
}

func getFullUri(pathPart string, accountSid string) string {
	return strings.Join([]string{BaseURL, APIVersion, "Accounts", accountSid, pathPart + ".json"}, "/")
}

// Convenience wrapper around MakeRequest
func (c *Client) GetResource(pathPart string, sid string, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest("GET", sidPart, nil, v)
}

func (c *Client) CreateResource(pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest("POST", pathPart, data, v)
}

func (c *Client) UpdateResource(pathPart string, sid string, data url.Values, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest("POST", sidPart, nil, v)
}

func (c *Client) ListResource(pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest("GET", pathPart, data, v)
}

// Make a request to the Twilio API.
func (c *Client) MakeRequest(method string, pathPart string, data url.Values, v interface{}) error {
	rb := new(strings.Reader)
	if data != nil && (method == "POST" || method == "PUT") {
		rb = strings.NewReader(data.Encode())
	}
	uri := getFullUri(pathPart, c.AccountSid)
	if method == "GET" && data != nil {
		uri = uri + "?" + data.Encode()
	}
	req, err := c.NewRequest(method, uri, rb)
	if err != nil {
		return err
	}
	return c.Do(req, &v)
}
