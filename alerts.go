package twilio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

const alertPathPart = "Alerts"

type AlertService struct {
	client *Client
}

type Alert struct {
	Sid        string `json:"sid"`
	AccountSid string `json:"account_sid"`
	// For Calls, AlertText is a series of key=value pairs separated by
	// ampersands
	AlertText        string          `json:"alert_text"`
	APIVersion       string          `json:"api_version"`
	DateCreated      TwilioTime      `json:"date_created"`
	DateGenerated    TwilioTime      `json:"date_generated"`
	DateUpdated      TwilioTime      `json:"date_updated"`
	ErrorCode        Code            `json:"error_code"`
	LogLevel         LogLevel        `json:"log_level"`
	MoreInfo         string          `json:"more_info"`
	RequestMethod    string          `json:"request_method"`
	RequestURL       string          `json:"request_url"`
	RequestVariables string          `json:"request_variables"`
	ResponseBody     string          `json:"response_body"`
	ResponseHeaders  Header          `json:"response_headers"`
	ResourceSid      string          `json:"resource_sid"`
	ServiceSid       json.RawMessage `json:"service_sid"`
	URL              string          `json:"url"`
}

type AlertPage struct {
	Meta   Meta     `json:"meta"`
	Alerts []*Alert `json:"alerts"`
}

func (a *AlertService) Get(ctx context.Context, sid string) (*Alert, error) {
	alert := new(Alert)
	err := a.client.GetResource(ctx, alertPathPart, sid, alert)
	return alert, err
}

func (a *AlertService) GetPage(ctx context.Context, data url.Values) (*AlertPage, error) {
	page := new(AlertPage)
	err := a.client.ListResource(ctx, alertPathPart, data, page)
	return page, err
}

type AlertPageIterator struct {
	p *PageIterator
}

func (a *AlertService) GetPageIterator(data url.Values) *AlertPageIterator {
	iter := NewPageIterator(a.client, data, alertPathPart)
	return &AlertPageIterator{
		p: iter,
	}
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (a *AlertPageIterator) Next(ctx context.Context) (*AlertPage, error) {
	ap := new(AlertPage)
	err := a.p.Next(ctx, ap)
	if err != nil {
		return nil, err
	}
	a.p.SetNextPageURI(ap.Meta.NextPageURL)
	return ap, nil
}

// Description tries as hard as possible to give you a one sentence description
// of this Alert, based on its contents. Description does not include a
// trailing period.
func (a *Alert) Description() string {
	vals, err := url.ParseQuery(a.AlertText)
	if err == nil {
		if msg := vals.Get("Msg"); msg != "" {
			return strings.TrimSpace(strings.TrimSuffix(msg, "."))
		}
	}
	if a.ErrorCode != 0 {
		switch a.ErrorCode {
		case CodeHTTPRetrievalFailure:
			s := "HTTP retrieval failure"
			if err == nil && vals.Get("httpResponse") != "" {
				s = fmt.Sprintf("%s: status code %s when fetching TwiML", s, vals.Get("httpResponse"))
			}
			return s
		default:
			if a.MoreInfo != "" {
				return fmt.Sprintf("Error %d: %s", a.ErrorCode, a.MoreInfo)
			}
			return fmt.Sprintf("Error %d", a.ErrorCode)
		}
	}
	if a.MoreInfo != "" {
		return "Unknown failure: " + a.MoreInfo
	}
	return "Unknown failure"
}
