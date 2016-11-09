package twilio

import (
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestGetMessagePrice(t *testing.T) {
	t.Parallel()
	client, server := getServer(messagePriceGB)
	defer server.Close()

	isoCountry := "GB"
	expectedCountryName := "United Kingdom"
	expectedPriceUnit := "USD"

	messagePrice, err := client.Pricing.Messaging.Countries.Get(context.Background(), isoCountry)
	if err != nil {
		t.Fatal(err)
	}
	if messagePrice == nil {
		t.Error("expected message price to be returned")
	}
	if messagePrice.Country != expectedCountryName {
		t.Errorf("Expected message price country to be %s, but got %s\n", expectedCountryName, messagePrice.Country)
	}
	if messagePrice.IsoCountry != isoCountry {
		t.Errorf("Expected message price iso country to be %s, but got %s\n", isoCountry, messagePrice.IsoCountry)
	}
	if messagePrice.PriceUnit != expectedPriceUnit {
		t.Errorf("Expected message price unit to be %s, but got %s\n", expectedPriceUnit, messagePrice.PriceUnit)
	}
	if messagePrice.InboundSmsPrices == nil {
		t.Error("Expected message price to contain InboundSmsPrices")
	}
	if messagePrice.OutboundSMSPrices == nil {
		t.Error("Expected message price to contain OutboundSMSPrices")
	}

	inboundPriceMap := make(map[string]bool)
	for _, inPrice := range messagePrice.InboundSmsPrices {
		inboundPriceMap[inPrice.NumberType] = true
	}
	// inboundPriceMap => map[local:true mobile:true shortcode:true]
	_, ok := inboundPriceMap["local"]
	if ok == false {
		t.Error("Expected inbound price to contain a price for local calls")
	}

	outboundCarrierPriceMap := make(map[string]bool)
	for _, outPrice := range messagePrice.OutboundSMSPrices {
		if outPrice.MCC == "" {
			t.Errorf("Each outbound SMS price should have a MCC, got %+v\n", outPrice)
		}
		if outPrice.MNC == "" {
			t.Errorf("Each outbound SMS price should have a MNC, got %+v\n", outPrice)
		}
		if outPrice.Carrier == "" {
			t.Errorf("Each outbound SMS price should have a Carrier, got %+v\n", outPrice)
		}
		if outPrice.Prices == nil {
			t.Errorf("Each outbound SMS price should have prices, got %+v\n", outPrice)
		}
		outboundCarrierPriceMap[outPrice.Carrier] = true
	}
	_, ok = outboundCarrierPriceMap["Other"]
	if ok == false {
		t.Error("Expected outbound price to contain a price for the carrier Other")
	}
}

func TestGetMessagingPricePage(t *testing.T) {
	t.Parallel()
	client, server := getServer(messagingCountriesPage)
	defer server.Close()

	data := url.Values{"PageSize": []string{"100"}}
	countriesPricePage, err := client.Pricing.Messaging.Countries.GetPage(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if len(countriesPricePage.Countries) == 0 {
		t.Error("expected to get a list of countries, got back 0")
	}
}
