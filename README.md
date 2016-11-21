# twilio-go

A client for accessing the Twilio API with several nice features:

- Easy-to-use helpers for purchasing phone numbers and sending MMS messages

- E.164 support, times that are parsed into a time.Time, and other smart types.

- Finer grained control over timeouts with a Context, and the library uses
  wall-clock HTTP timeouts, not socket timeouts.

- Easy debugging network traffic by setting DEBUG_HTTP_TRAFFIC=true in your
  environment.

- Easily find calls and messages that occurred between a particular
set of `time.Time`s, down to the nanosecond, with GetCallsInRange /
GetMessagesInRange.

- It's clear when the library will make a network request, there are no
unexpected latency spikes when paging from one resource to the next.

- Uses threads to fetch resources concurrently; for example, has methods to
fetch all Media for a Message concurrently.

- Usable, [one sentence descriptions of Alerts][alert-descriptions].

[alert-descriptions]: https://godoc.org/github.com/saintpete/twilio-go#Alert.Description

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
godoc.org](https://godoc.org/github.com/saintpete/twilio-go).

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
- Pricing

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
    StatusCode: 404
}
```

Not all errors will be a `rest.Error` however - HTTP timeouts, canceled
context.Contexts, and JSON parse errors (HTML error pages, bad gateway
responses from proxies) may also be returned as plain Go errors.

### Twiml Generation

There are no plans to support Twiml generation in this library. It may be
more readable and maintainable to manually write the XML involved in a Twiml
response.

### API Problems this Library Solves For You

- Media URL's are returned over HTTP. twilio-go rewrites these URL's to be
  HTTPS before returning them to you.

- A subset of Notifications returned code 4107, which doesn't exist. These
  notifications should have error code 14107. We rewrite the error code
  internally before returning it to you.

- The only provided API for filtering calls or messages by date grabs all
messages for an entire day, and the day ranges are only available for UTC. Use
GetCallsInRange or GetMessagesInRange to do timezone-aware, finer-grained date
filtering.

### Errata

You can get Alerts for a given Call or MMS by passing `ResourceSid=CA123` as
a filter to Alerts.GetPage. This functionality is not documented in the API.

[zero-results]: https://github.com/saintpete/twilio-go/issues/8
