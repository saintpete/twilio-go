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

The API is experimental and subject to change at any point.
