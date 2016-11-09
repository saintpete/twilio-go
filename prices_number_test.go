package twilio

import (
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestGetPhoneNumberPriceGB(t *testing.T) {
	t.Parallel()
	client, server := getServer(phoneNumberPriceGB)
	defer server.Close()

	isoCountry := "GB"
	expectedCountryName := "United Kingdom"
	expectedPriceUnit := "USD"

	numPrice, err := client.Pricing.PhoneNumbers.Countries.Get(context.Background(), isoCountry)
	if err != nil {
		t.Fatal(err)
	}
	if numPrice == nil {
		t.Error("expected voice price to be returned")
	}
	if numPrice.Country != expectedCountryName {
		t.Errorf("Expected voice price country to be %s, but got %s\n", expectedCountryName, numPrice.Country)
	}
	if numPrice.IsoCountry != isoCountry {
		t.Errorf("Expected voice price iso country to be %s, but got %s\n", isoCountry, numPrice.IsoCountry)
	}
	if numPrice.PriceUnit != expectedPriceUnit {
		t.Errorf("Expected voice price unit to be %s, but got %s\n", expectedPriceUnit, numPrice.PriceUnit)
	}
	if numPrice.PhoneNumberPrices == nil {
		t.Error("Expected voice price to contain PhoneNumberPrices")
	}

	numTypePriceMap := make(map[string]bool)
	for _, price := range numPrice.PhoneNumberPrices {
		numTypePriceMap[price.NumberType] = true
	}
	// numTypePriceMap => map[mobile:true national:true toll free:true local:true]
	_, ok := numTypePriceMap["local"]
	if ok == false {
		t.Error("Expected number price to contain a price for a local number")
	}
}

func TestGetPhoneNumbersPricePage(t *testing.T) {
	t.Parallel()
	client, server := getServer(phoneNumberCountriesPage)
	defer server.Close()

	data := url.Values{"PageSize": []string{"100"}}
	countriesPricePage, err := client.Pricing.PhoneNumbers.Countries.GetPage(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if len(countriesPricePage.Countries) == 0 {
		t.Error("expected to get a list of countries, got back 0")
	}
}
