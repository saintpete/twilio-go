package twilio_test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kevinburke/rest"
	twilio "github.com/kevinburke/twilio-go"
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
