package twilio

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	types "github.com/kevinburke/go-types"
	"github.com/kevinburke/rest"
)

const Version = "0.13"
const userAgent = "twilio-go/" + Version

var BaseURL = "https://api.twilio.com"

const APIVersion = "2010-04-01"

type Client struct {
	*rest.Client

	AccountSid string
	AuthToken  string

	Media    *MediaService
	Messages *MessageService
}

type Page struct {
	FirstPageURI string           `json:"first_page_uri"`
	Start        uint             `json:"start"`
	End          uint             `json:"end"`
	NumPages     uint             `json:"num_pages"`
	Total        uint             `json:"total"`
	NextPageURI  types.NullString `json:"next_page_uri"`
	PageSize     uint             `json:"page_size"`
}

// NoMoreResults is returned if you reach the end of the result set while
// paging through resources.
var NoMoreResults = errors.New("twilio: No more results")

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
	c.Media = &MediaService{client: c}
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
