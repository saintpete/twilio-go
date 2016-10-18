package twilio

import (
	"fmt"
	"testing"
)

func TestGetCall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	call, err := envClient.Calls.Get("CAa98f7bbc9bc4980a44b128ca4884ca73")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", call)
}
