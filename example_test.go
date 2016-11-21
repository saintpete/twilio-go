package twilio_test

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/kevinburke/rest"
	twilio "github.com/saintpete/twilio-go"
	"golang.org/x/net/context"
)

var callURL, _ = url.Parse("https://kev.inburke.com/zombo/zombocom.mp3")

func Example() {
	client := twilio.NewClient("AC123", "123", nil)

	// Send a SMS
	msg, _ := client.Messages.SendMessage("+14105551234", "+14105556789", "Sent via go :) ✓", nil)
	fmt.Println(msg.Sid, msg.FriendlyPrice())

	// Make a call
	call, _ := client.Calls.MakeCall("+14105551234", "+14105556789", callURL)
	fmt.Println(call.Sid, call.FriendlyPrice())

	_, err := client.IncomingNumbers.BuyNumber("+1badnumber")
	// Twilio API errors are converted to rest.Error types
	if err != nil {
		restErr, ok := err.(*rest.Error)
		if ok {
			fmt.Println(restErr.Title)
			fmt.Println(restErr.Type)
		}
	}

	// Find all calls from a number
	data := url.Values{"From": []string{"+14105551234"}}
	iterator := client.Calls.GetPageIterator(data)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for {
		page, err := iterator.Next(ctx)
		if err == twilio.NoMoreResults {
			break
		}
		for _, call := range page.Calls {
			fmt.Println(call.Sid, call.To)
		}
	}
}

func ExampleCallService_GetCallsInRange() {
	// Get all calls between 10:34:00 Oct 26 and 19:25:59 Oct 27, NYC time.
	nyc, _ := time.LoadLocation("America/New_York")
	start := time.Date(2016, 10, 26, 22, 34, 00, 00, nyc)
	end := time.Date(2016, 10, 27, 19, 25, 59, 00, nyc)

	client := twilio.NewClient("AC123", "123", nil)
	iter := client.Calls.GetCallsInRange(start, end, url.Values{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		page, err := iter.Next(ctx)
		if err == twilio.NoMoreResults {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i, call := range page.Calls {
			fmt.Printf("%d: %s (%s)", i, call.Sid, call.DateCreated.Time)
		}
	}
}
