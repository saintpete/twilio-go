package twilio

import (
	"golang.org/x/net/context"
	"net/url"
)

const availablePhoneNumbersPath = "AvailablePhoneNumbers"

type AvailablePhoneNumberBase struct {
	client   *Client
	pathPart string
}

type AvailablePhoneNumberService struct {
	Local    *AvailablePhoneNumberBase
	Mobile   *AvailablePhoneNumberBase
	TollFree *AvailablePhoneNumberBase
}

// The subresources of the AvailablePhoneNumbers resource let you search for local, toll-free and
// mobile phone numbers that are available for you to purchase.
// See https://www.twilio.com/docs/api/rest/available-phone-numbers for details
type AvailablePhoneNumber struct {
	FriendlyName        string            `json:"friendly_name"`
	PhoneNumber         PhoneNumber       `json:"phone_number"`
	Lata                string            `json:"lata"`
	RateCenter          string            `json:"rate_center"`
	Latitude            string            `json:"latitude"`
	Longitude           string            `json:"longitude"`
	Region              string            `json:"region"`
	PostalCode          string            `json:"postal_code"`
	ISOCountry          string            `json:"iso_country"`
	Capabilities        *NumberCapability `json:"capabilities"`
	AddressRequirements string            `json:"address_requirements"`
	Beta                bool              `json:"beta"`
}

type AvailablePhoneNumberPage struct {
	Uri     string                  `json:"uri"`
	Numbers []*AvailablePhoneNumber `json:"available_phone_numbers"`
}

// https://www.twilio.com/docs/api/rest/available-phone-numbers#local
// https://www.twilio.com/docs/api/rest/available-phone-numbers#toll-free
// https://www.twilio.com/docs/api/rest/available-phone-numbers#mobile
func (s *AvailablePhoneNumberBase) GetPage(ctx context.Context, isoCountry string, filters url.Values) (*AvailablePhoneNumberPage, error) {
	sr := new(AvailablePhoneNumberPage)
	path := availablePhoneNumbersPath + "/" + isoCountry + "/" + s.pathPart
	err := s.client.ListResource(ctx, path, filters, sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}
