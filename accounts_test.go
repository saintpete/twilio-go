package twilio

import (
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestAccountGet(t *testing.T) {
	t.Parallel()
	client, server := getServer(accountInstance)
	defer server.Close()
	sid := "AC58f1e8f2b1c6b88ca90a012a4be0c279"
	acct, err := client.Accounts.Get(context.Background(), sid)
	if err != nil {
		t.Fatal(err)
	}
	if acct.Sid != sid {
		t.Errorf("wrong sid")
	}
	if acct.Type != "Full" {
		t.Errorf("wrong type")
	}
}

func TestAccountGetPage(t *testing.T) {
	t.Parallel()
	client, server := getServer(accountList)
	defer server.Close()
	data := url.Values{}
	data.Set("PageSize", "2")
	page, err := client.Accounts.GetPage(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Accounts) != 2 {
		t.Errorf("expected accounts len to be 2")
	}
	if page.Accounts[0].Sid != "AC0cd9be8fd5e6e4fa0a04f50ac1caca4e" {
		t.Errorf("wrong sid")
	}
}
