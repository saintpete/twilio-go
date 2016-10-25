# twilio-go

A client for accessing the Twilio API with several nice features:

- Can fetch all Media for a Message concurrently

- Easy-to-use helpers for purchasing phone numbers and sending MMS messages

- E.164 support and other smart types.

- Wall-clock HTTP timeouts, vs. socket timeouts

- Easy debugging network traffic by setting DEBUG_HTTP_TRAFFIC=true in your
  environment.

Here are some example use cases:

```go
const sid = "AC123"
const token = "456bef"

client := twilio.NewClient(sid, token, nil)

// Send a message
msg, err := client.Messages.SendMessage("+14105551234", "+14105556789", "Sent via go :) ✓", nil)

// Buy a number
number, err := client.IncomingNumbers.BuyNumber("+14105551234")

// Get all calls from a number
data := new(url.Values)
data.Set("From", "+14105551234")
callPage, err := client.Calls.GetPage(data)

// Iterate over calls
iterator := client.Calls.GetPageIterator(url.Values{})
for {
    page, err := iterator.Next()
    if err == twilio.NoMoreResults {
        break
    }
    fmt.Println("start", page.Start)
}
```

A [complete documentation reference can be found at
godoc.org](https://godoc.org/github.com/kevinburke/twilio-go).

The API is experimental, but unlikely to change, and currently only covers
these resources:

- Calls
- Messages
- Incoming Phone Numbers
- Recordings
- Media

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

Not all errors will be a `rest.Error` however - HTTP timeouts and JSON parse
errors (HTML error pages, bad gateway responses from proxies) can also be
returned.
