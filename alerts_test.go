package twilio

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGetAlert(t *testing.T) {
	t.Parallel()
	client, s := getServer(alertInstanceResponse)
	defer s.Close()
	sid := "NO00ed1fb4aa449be2434d54ec8e492349"
	alert, err := client.Monitor.Alerts.Get(context.Background(), sid)
	if err != nil {
		t.Fatal(err)
	}
	if alert.Sid != sid {
		t.Errorf("expected Sid to be %s, got %s", sid, alert.Sid)
	}
}

func TestGetAlertPage(t *testing.T) {
	t.Parallel()
	client, s := getServer(alertListResponse)
	defer s.Close()
	data := url.Values{"PageSize": []string{"2"}}
	page, err := client.Monitor.Alerts.GetPage(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Alerts) != 2 {
		t.Errorf("expected 2 alerts, got %d", len(page.Alerts))
	}
	if page.Meta.Key != "alerts" {
		t.Errorf("expected Key to be 'alerts', got %s", page.Meta.Key)
	}
	if page.Meta.PageSize != 2 {
		t.Errorf("expected PageSize to be 2, got %d", page.Meta.PageSize)
	}
	if page.Meta.Page != 0 {
		t.Errorf("expected Page to be 0, got %d", page.Meta.Page)
	}
	if page.Meta.PreviousPageURL.Valid != false {
		t.Errorf("expected previousPage.Valid to be false, got true")
	}
	if page.Alerts[0].LogLevel != LogLevelError {
		t.Errorf("expected LogLevel to be Error, got %s", page.Alerts[0].LogLevel)
	}
	if page.Alerts[0].RequestMethod != "POST" {
		t.Errorf("expected RequestMethod to be 'POST', got %s", page.Alerts[0].RequestMethod)
	}
	if page.Alerts[0].ErrorCode != CodeHTTPRetrievalFailure {
		t.Errorf("expected ErrorCode to be '11200', got %d", page.Alerts[0].ErrorCode)
	}
}

func TestGetAlertIterator(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	// TODO make a better assertion here or a server that can return different
	// responses
	iter := envClient.Monitor.Alerts.GetPageIterator(url.Values{"PageSize": []string{"35"}})
	count := uint(0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for {
		page, err := iter.Next(ctx)
		if err == NoMoreResults {
			break
		}
		if err != nil {
			t.Fatal(err)
			break
		}
		if page.Meta.Page < count {
			t.Fatal("too small page number")
			return
		}
		count++
		if count > 15 {
			fmt.Println("count > 15")
			t.Fail()
			break
		}
	}
	if count < 10 {
		t.Errorf("Too small of a count - expected at least 10, got %d", count)
	}
}
