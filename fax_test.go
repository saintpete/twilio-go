package twilio

import (
	"context"
	"net/url"
	"testing"
)

func TestGetFax(t *testing.T) {
	t.Parallel()
	client, server := getServer(faxGetResponse)
	defer server.Close()
	fax, err := client.Fax.Faxes.Get(context.Background(), "FXeb76f282888a074547beba3516552174")
	if err != nil {
		t.Fatal(err)
	}
	if fax.Quality != "fine" {
		t.Errorf("quality is incorrect")
	}
	if fax.FriendlyPrice() != "$0.014" {
		t.Errorf("FriendlyPrice(): got %s, want $0.014", fax.FriendlyPrice())
	}
	if fax.NumPages != 1 {
		t.Errorf("wrong NumPages %d", fax.NumPages)
	}
	if fax.Status != StatusDelivered {
		t.Errorf("wrong status: %s", fax.Status)
	}
}

func TestGetFaxPage(t *testing.T) {
	t.Parallel()
	client, server := getServer(faxGetPageResponse)
	defer server.Close()
	page, err := client.Fax.Faxes.GetPage(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Faxes) == 0 {
		t.Errorf("too few faxes returned")
	}
}

func TestSendFax(t *testing.T) {
	t.Parallel()
	client, server := getServer(faxCreateResponse)
	defer server.Close()
	data := url.Values{
		"To":             []string{"+18326327228"},
		"From":           []string{"+19252717005"},
		"MediaUrl":       []string{"https://kb.ngrok.io/gopher.pdf"},
		"Quality":        []string{"fine"},
		"StatusCallback": []string{"https://kb.ngrok.io/fax-callback"},
	}
	fax, err := client.Fax.Faxes.Create(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if fax.Quality != "fine" {
		t.Errorf("quality is incorrect")
	}
}
