package twilio

import (
	"net/url"

	"golang.org/x/net/context"
)

// We need to build this relative to the root, but users can override the
// APIVersion, so give them a chance to before our init runs.
var accountPathPart string

func init() {
	accountPathPart = "/" + APIVersion + "/Accounts"
}

type Account struct {
	Sid             string            `json:"sid"`
	FriendlyName    string            `json:"friendly_name"`
	Type            string            `json:"type"`
	AuthToken       string            `json:"auth_token"`
	OwnerAccountSid string            `json:"owner_account_sid"`
	DateCreated     TwilioTime        `json:"date_created"`
	DateUpdated     TwilioTime        `json:"date_updated"`
	Status          Status            `json:"status"`
	SubresourceURIs map[string]string `json:"subresource_uris"`
	URI             string            `json:"uri"`
}

type AccountPage struct {
	Page
	Accounts []*Account `json:"accounts"`
}

type AccountService struct {
	client *Client
}

func (a *AccountService) Get(ctx context.Context, sid string) (*Account, error) {
	acct := new(Account)
	// hack because this is not a resource off of the account sid
	sidJSON := sid + ".json"
	err := a.client.GetResource(ctx, accountPathPart, sidJSON, acct)
	return acct, err
}

func (a *AccountService) GetPage(ctx context.Context, data url.Values) (*AccountPage, error) {
	ap := new(AccountPage)
	err := a.client.ListResource(ctx, accountPathPart+".json", data, ap)
	return ap, err
}
