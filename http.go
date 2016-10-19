package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kevinburke/rest"
)

const Version = "0.24"
const userAgent = "twilio-go/" + Version

var BaseURL = "https://api.twilio.com"

const APIVersion = "2010-04-01"

type Client struct {
	*rest.Client

	AccountSid string
	AuthToken  string

	Calls      *CallService
	Media      *MediaService
	Messages   *MessageService
	Recordings *RecordingService
}

const defaultTimeout = 30*time.Second + 500*time.Millisecond

// NewClient creates a Client for interacting with the Twilio API.
func NewClient(accountSid string, authToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	restClient := rest.NewClient(accountSid, authToken, fmt.Sprintf("%s/%s", BaseURL, APIVersion))
	restClient.Client = httpClient
	restClient.UploadType = rest.FormURLEncoded

	c := &Client{Client: restClient, AccountSid: accountSid, AuthToken: authToken}
	c.Calls = &CallService{client: c}
	c.Media = &MediaService{client: c}
	c.Messages = &MessageService{client: c}
	c.Recordings = &RecordingService{client: c}
	return c
}

func (c *Client) FullPath(pathPart string) string {
	return "/" + strings.Join([]string{"Accounts", c.AccountSid, pathPart + ".json"}, "/")
}

// GetResource retrieves an instance resource with the given path part (e.g.
// "/Messages") and sid (e.g. "SM123").
func (c *Client) GetResource(pathPart string, sid string, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest("GET", sidPart, nil, v)
}

// CreateResource makes a POST request to the given resource.
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

func (c *Client) GetNextPage(fullUri string, v interface{}) error {
	return c.MakeRequest("GET", fullUri, nil, v)
}

// Make a request to the Twilio API.
func (c *Client) MakeRequest(method string, pathPart string, data url.Values, v interface{}) error {
	if strings.HasPrefix(pathPart, "/"+APIVersion) {
		pathPart = pathPart[len("/"+APIVersion):]
	} else {
		pathPart = c.FullPath(pathPart)
	}
	rb := new(strings.Reader)
	if data != nil && (method == "POST" || method == "PUT") {
		rb = strings.NewReader(data.Encode())
	}
	if method == "GET" && data != nil {
		pathPart = pathPart + "?" + data.Encode()
	}
	req, err := c.NewRequest(method, pathPart, rb)
	if err != nil {
		return err
	}
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", userAgent+" "+ua)
	}
	return c.Do(req, &v)
}
