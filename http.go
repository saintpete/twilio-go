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
const Version = "0.55"
const userAgent = "twilio-go/" + Version

// The base URL serving the API. Override this for testing.
var BaseURL = "https://api.twilio.com"

// The base URL for Twilio Monitor.
var MonitorBaseURL = "https://monitor.twilio.com"

// Version of the Twilio Monitor API.
const MonitorVersion = "v1"

// The base URL for Twilio Pricing.
var PricingBaseURL = "https://pricing.twilio.com"

// Version of the Twilio Pricing API.
const PricingVersion = "v1"

// The APIVersion to use. Your mileage may vary using other values for the
// APIVersion; the resource representations may not match.
const APIVersion = "2010-04-01"

type Client struct {
	*rest.Client
	Monitor *Client
	Pricing *Client

	// FullPath takes a path part (e.g. "Messages") and
	// returns the full API path, including the version (e.g.
	// "/2010-04-01/Accounts/AC123/Messages").
	FullPath func(pathPart string) string
	// The API version.
	APIVersion string

	AccountSid string
	AuthToken  string

	// The API Client uses these resources
	Accounts          *AccountService
	Applications      *ApplicationService
	Calls             *CallService
	Conferences       *ConferenceService
	IncomingNumbers   *IncomingNumberService
	Keys              *KeyService
	Media             *MediaService
	Messages          *MessageService
	OutgoingCallerIDs *OutgoingCallerIDService
	Queues            *QueueService
	Recordings        *RecordingService
	Transcriptions    *TranscriptionService

	// NewMonitorClient initializes these services
	Alerts *AlertService

	// NewPricingClient initializes these services
	Voice        *VoicePriceService
	Messaging    *MessagingPriceService
	PhoneNumbers *PhoneNumberPriceService
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
		Status: resp.StatusCode,
	}
}

func NewMonitorClient(accountSid string, authToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	restClient := rest.NewClient(accountSid, authToken, MonitorBaseURL)
	c := &Client{Client: restClient, AccountSid: accountSid, AuthToken: authToken}
	c.FullPath = func(pathPart string) string {
		return "/" + c.APIVersion + "/" + pathPart
	}
	c.APIVersion = MonitorVersion
	c.Alerts = &AlertService{client: c}
	return c
}

// returns a new Client to use the pricing API
func NewPricingClient(accountSid string, authToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	restClient := rest.NewClient(accountSid, authToken, PricingBaseURL)
	c := &Client{Client: restClient, AccountSid: accountSid, AuthToken: authToken}
	c.APIVersion = PricingVersion
	c.FullPath = func(pathPart string) string {
		return "/" + c.APIVersion + "/" + pathPart
	}
	c.Voice = &VoicePriceService{
		Countries: &CountryVoicePriceService{client: c},
		Numbers:   &NumberVoicePriceService{client: c},
	}
	c.Messaging = &MessagingPriceService{
		Countries: &CountryMessagingPriceService{client: c},
	}
	c.PhoneNumbers = &PhoneNumberPriceService{
		Countries: &CountryPhoneNumberPriceService{client: c},
	}
	return c
}

// NewClient creates a Client for interacting with the Twilio API. This is the
// main entrypoint for API interactions; view the methods on the subresources
// for more information.
func NewClient(accountSid string, authToken string, httpClient *http.Client) *Client {

	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	restClient := rest.NewClient(accountSid, authToken, BaseURL)
	restClient.Client = httpClient
	restClient.UploadType = rest.FormURLEncoded
	restClient.ErrorParser = parseTwilioError

	c := &Client{Client: restClient, AccountSid: accountSid, AuthToken: authToken}
	c.APIVersion = APIVersion

	c.FullPath = func(pathPart string) string {
		return "/" + strings.Join([]string{c.APIVersion, "Accounts", c.AccountSid, pathPart + ".json"}, "/")
	}
	c.Monitor = NewMonitorClient(accountSid, authToken, httpClient)
	c.Pricing = NewPricingClient(accountSid, authToken, httpClient)

	c.Accounts = &AccountService{client: c}
	c.Applications = &ApplicationService{client: c}
	c.Calls = &CallService{client: c}
	c.Conferences = &ConferenceService{client: c}
	c.Keys = &KeyService{client: c}
	c.Media = &MediaService{client: c}
	c.Messages = &MessageService{client: c}
	c.OutgoingCallerIDs = &OutgoingCallerIDService{client: c}
	c.Queues = &QueueService{client: c}
	c.Recordings = &RecordingService{client: c}
	c.Transcriptions = &TranscriptionService{client: c}

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

// RequestOnBehalfOf will make all future client requests using the same
// Account Sid and Auth Token for Basic Auth, but will use the provided
// subaccountSid in the URL. Use this to make requests on behalf of a
// subaccount, using the parent account's credentials.
//
// RequestOnBehalfOf is *not* thread safe, and modifies the Client's behavior
// for all requests going forward.
//
// RequestOnBehalfOf should only be used with api.twilio.com, not (for example)
// Twilio Monitor.
//
// To authenticate using a subaccount sid / auth token, create a new Client
// using that account's credentials.
func (c *Client) RequestOnBehalfOf(subaccountSid string) {
	c.FullPath = func(pathPart string) string {
		return "/" + strings.Join([]string{c.APIVersion, "Accounts", subaccountSid, pathPart + ".json"}, "/")
	}
}

// GetResource retrieves an instance resource with the given path part (e.g.
// "/Messages") and sid (e.g. "MM123").
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
	return c.MakeRequest(ctx, "POST", sidPart, data, v)
}

func (c *Client) DeleteResource(ctx context.Context, pathPart string, sid string) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	err := c.MakeRequest(ctx, "DELETE", sidPart, nil, nil)
	if err == nil {
		return nil
	}
	rerr, ok := err.(*rest.Error)
	if ok && rerr.StatusCode == http.StatusNotFound {
		return nil
	}
	return err
}

func (c *Client) ListResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "GET", pathPart, data, v)
}

// GetNextPage fetches the Page at fullUri and decodes it into v. fullUri
// should be a next_page_uri returned in the response to a paging request, and
// should be the full path, eg "/2010-04-01/.../Messages?Page=1&PageToken=..."
func (c *Client) GetNextPage(ctx context.Context, fullUri string, v interface{}) error {
	// for monitor etc.
	if strings.HasPrefix(fullUri, c.Base) {
		fullUri = fullUri[len(c.Base):]
	}
	return c.MakeRequest(ctx, "GET", fullUri, nil, v)
}

// Make a request to the Twilio API.
func (c *Client) MakeRequest(ctx context.Context, method string, pathPart string, data url.Values, v interface{}) error {
	if !strings.HasPrefix(pathPart, "/"+c.APIVersion) {
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
