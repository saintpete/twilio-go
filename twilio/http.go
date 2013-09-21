package twilio

import (
	"encoding/json"
	//"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const BaseUrl = "https://api.twilio.com"
const Version = "2010-04-01"

type Client struct {
	http.Client

	AccountSid string
	AuthToken  string

	Messages *MessageService
}

func CreateClient(accountSid string, authToken string, httpClient *http.Client) *Client {
	tr := &http.Transport{
		ResponseHeaderTimeout: time.Duration(3050) * time.Millisecond,
	}

	if httpClient == nil {
		httpClient = &http.Client{Transport: tr}
	}

	c := &Client{AccountSid: accountSid, AuthToken: authToken}
	c.Client = *httpClient

	c.Messages = &MessageService{client: c}
	return c
}

func getFullUri(pathPart string, accountSid string) string {
	return strings.Join([]string{BaseUrl, Version, "Accounts", accountSid, pathPart + ".json"}, "/")
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

// Make a request to the Twilio API.
func (c *Client) MakeRequest(method string, pathPart string, data url.Values, v interface{}) error {
	req, err := c.CreateRequest(method, pathPart, data)
	if err != nil {
		glog.Errorf("Error creating request", err)
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		glog.Errorf("Error making request", err)
		return err
	}
	defer resp.Body.Close()

	// XXX investigate whether this could overflow with a large response body,
	// it appears so from reading the ioutil source.
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &v)
	return nil
}

// Initializes the http request.
func (c *Client) CreateRequest(method string, pathPart string, data url.Values) (*http.Request, error) {
	var rb strings.Reader
	if data != nil && (method == "POST" || method == "PUT") {
		rb = *strings.NewReader(data.Encode())
	}
	uri := getFullUri(pathPart, c.AccountSid)
	req, err := http.NewRequest(method, uri, &rb)
	if err != nil {
		glog.Errorf("Couldn't parse api.twilio.com as a URL. This should not happen", err)
		return nil, err
	}
	req.SetBasicAuth(c.AccountSid, c.AuthToken)
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Add("Accept", "application/json")
	// XXX add system platform information
	req.Header.Add("User-Agent", "twilio-go/0.0.1")
	if data != nil && (method == "POST" || method == "PUT") {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}
