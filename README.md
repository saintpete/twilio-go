# twilio-go

Usage

```go
const sid = "AC123"
const token = "456bef"

client := twilio.CreateClient(sid, token, nil)
msg, err := client.Messages.SendMessage("+19177460767", "+19253240627", "Sent via go :) ✓", nil)
```

Or pass the values yourself:

```go
params := url.Values{
    "From": {"+19177460767"},
    "To":   {"+19253240627"},
    "Body": {"Sent via go :) ✓"},
}
msg, err := client.Messages.Create(params)
```

The API is experimental, doesn't cover all resources, and subject to change
until 1.0.

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
