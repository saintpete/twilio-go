# twilio-go

A client for accessing the Twilio API with several nice features:

- Easy-to-use helpers for purchasing phone numbers and sending MMS messages

- E.164 support and other smart types.

- Finer grained control over timeouts with a Context, and the library uses
  wall-clock HTTP timeouts, not socket timeouts.

- Easy debugging network traffic by setting DEBUG_HTTP_TRAFFIC=true in your
  environment.

- Clarity; it's clear when the library will make a network request, there's no
  unexpected giant latency spikes when paging through resources.

- Uses threads to fetch resources concurrently; for example, has methods to
fetch all Media for a Message concurrently.

Here are some example use cases:

```go
const sid = "AC123"
const token = "456bef"

client := twilio.NewClient(sid, token, nil)

// Send a message
msg, err := client.Messages.SendMessage("+14105551234", "+14105556789", "Sent via go :) ✓", nil)

// Start a phone call
call, err := client.Calls.MakeCall("+14105551234", "+14105556789",
        "https://kev.inburke.com/zombo/zombocom.mp3")

// Buy a number
number, err := client.IncomingNumbers.BuyNumber("+14105551234")

// Get all calls from a number
data := url.Values{}
data.Set("From", "+14105551234")
callPage, err := client.Calls.GetPage(context.TODO(), data)

// Iterate over calls
iterator := client.Calls.GetPageIterator(url.Values{})
for {
    page, err := iterator.Next(context.TODO())
    if err == twilio.NoMoreResults {
        break
    }
    fmt.Println("start", page.Start)
}
```

A [complete documentation reference can be found at
godoc.org](https://godoc.org/github.com/kevinburke/twilio-go).

The API is open to, but unlikely to change, and currently only covers
these resources:

- Alerts
- Applications
- Calls
- Conferences
- Incoming Phone Numbers
- Keys
- Messages
- Media
- Outgoing Caller ID's
- Queues
- Recordings
- Transcriptions
- Access Tokens for IPMessaging, Video and Programmable Voice SDK

### Error Parsing

If the twilio-go client gets an error from
the Twilio API, we attempt to convert it to a
[`rest.Error`](https://godoc.org/github.com/kevinburke/rest#Error) before
returning. Here's an example 404.

```
&rest.Error{
    Title: "The requested resource ... was not found",
    ID: "20404",
    Detail: "",
    Instance: "",
    Type: "https://www.twilio.com/docs/errors/20404",
    StatusCode:404
}
```

Not all errors will be a `rest.Error` however - HTTP timeouts, canceled
context.Contexts, and JSON parse errors (HTML error pages, bad gateway
responses from proxies) may also be returned as plain Go errors.
