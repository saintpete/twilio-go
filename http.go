package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kevinburke/rest"
	"golang.org/x/net/context"
)

// The twilio-go version. Run "make release" to bump this number.
const Version = "0.31"
const userAgent = "twilio-go/" + Version

// The base URL serving the API. Override this for testing.
var BaseURL = "https://api.twilio.com"

// The APIVersion to use. Your mileage may vary using other values for the
// APIVersion; the resource representations may not match.
const APIVersion = "2010-04-01"

type Client struct {
	*rest.Client

	AccountSid string
	AuthToken  string

	Calls           *CallService
	Conferences     *ConferenceService
	IncomingNumbers *IncomingNumberService
	Media           *MediaService
	Messages        *MessageService
	Recordings      *RecordingService
}

const defaultTimeout = 30*time.Second + 500*time.Millisecond

// An error returned by the Twilio API. We don't want to expose this - let's
// try to standardize on the fields in the HTTP problem spec instead.
type twilioError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	// This will be ignored in favor of the actual HTTP status code
	Status int `json:"status"`
}

func parseTwilioError(resp *http.Response) error {
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	rerr := new(twilioError)
	err = json.Unmarshal(resBody, rerr)
	if err != nil {
		return fmt.Errorf("invalid response body: %s", string(resBody))
	}
	if rerr.Message == "" {
		return fmt.Errorf("invalid response body: %s", string(resBody))
	}
	return &rest.Error{
		Title:      rerr.Message,
		Type:       rerr.MoreInfo,
		ID:         strconv.FormatInt(int64(rerr.Code), 10),
		StatusCode: resp.StatusCode,
	}
}

// NewClient creates a Client for interacting with the Twilio API. This is the
// main entrypoint for API interactions; view the methods on the subresources
// for more information.
func NewClient(accountSid string, authToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	restClient := rest.NewClient(accountSid, authToken, fmt.Sprintf("%s/%s", BaseURL, APIVersion))
	restClient.Client = httpClient
	restClient.UploadType = rest.FormURLEncoded
	restClient.ErrorParser = parseTwilioError

	c := &Client{Client: restClient, AccountSid: accountSid, AuthToken: authToken}
	c.Calls = &CallService{client: c}
	c.Conferences = &ConferenceService{client: c}
	c.Media = &MediaService{client: c}
	c.Messages = &MessageService{client: c}
	c.Recordings = &RecordingService{client: c}

	c.IncomingNumbers = &IncomingNumberService{
		NumberPurchasingService: &NumberPurchasingService{
			client:   c,
			pathPart: "",
		},
		client: c,
		Local: &NumberPurchasingService{
			client:   c,
			pathPart: "Local",
		},
		TollFree: &NumberPurchasingService{
			client:   c,
			pathPart: "TollFree",
		},
	}
	return c
}

func (c *Client) FullPath(pathPart string) string {
	return "/" + strings.Join([]string{"Accounts", c.AccountSid, pathPart + ".json"}, "/")
}

// GetResource retrieves an instance resource with the given path part (e.g.
// "/Messages") and sid (e.g. "SM123").
func (c *Client) GetResource(ctx context.Context, pathPart string, sid string, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest(ctx, "GET", sidPart, nil, v)
}

// CreateResource makes a POST request to the given resource.
func (c *Client) CreateResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "POST", pathPart, data, v)
}

func (c *Client) UpdateResource(ctx context.Context, pathPart string, sid string, data url.Values, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest(ctx, "POST", sidPart, nil, v)
}

func (c *Client) ListResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "GET", pathPart, data, v)
}

func (c *Client) GetNextPage(ctx context.Context, fullUri string, v interface{}) error {
	return c.MakeRequest(ctx, "GET", fullUri, nil, v)
}

// Make a request to the Twilio API.
func (c *Client) MakeRequest(ctx context.Context, method string, pathPart string, data url.Values, v interface{}) error {
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
	req = withContext(req, ctx)
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", userAgent+" "+ua)
	}
	return c.Do(req, &v)
}
