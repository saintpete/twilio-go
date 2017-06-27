// Binary report-data-usage reports information about how much data your sim cards have used.
//
// Example report format:
//
//     $ report-data-usage
//     iPhone v13
//     ----------
//     2017-06-19: 0
//     2017-06-20: 0
//     2017-06-21: 11.281MB
//     2017-06-22: 93.341MB
//     2017-06-23: 94.422MB
//     2017-06-24: 159.461MB
//     2017-06-25: 50.062MB
//     total (last 7 days): 408.568MB
//
// Your Twilio credentials are loaded from the TWILIO_ACCOUNT_SID and
// TWILIO_AUTH_TOKEN environment variables. There are several flags:
//
//     --days int
//         Change the number of days to report usage for (default 7)
//     --location string
//         Use a different timezone for day boundaries (example "America/Los_Angeles")
//     --sim string
//         Only fetch usage for this sim
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	types "github.com/kevinburke/go-types"
	twilio "github.com/kevinburke/twilio-go"
	"golang.org/x/sync/errgroup"
)

var duration = flag.Uint("days", 7, "Number of days to get usage for")
var location = flag.String("location", "", "Timezone to use (defaults to system location/TZ env var)")
var sim = flag.String("sim", "", "Load data about a specific Sim/unique name only")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Print information about how much data you used on each day.
`)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	var loc *time.Location
	if *location == "" {
		loc = time.Local
	} else {
		var err error
		loc, err = time.LoadLocation(*location)
		if err != nil {
			log.Fatal(err)
		}
	}
	c := twilio.NewClient(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), nil)
	now := time.Now().Add(24 * time.Hour).In(loc)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	usage := make(map[string][]types.Bits)
	iter := c.Wireless.Sims.GetPageIterator(nil)
	sims := make([]string, 0)
	if *sim != "" {
		sims = []string{*sim}
	} else {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 31*time.Second)
			page, err := iter.Next(ctx)
			cancel()
			if err == twilio.NoMoreResults {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < len(page.Sims); i++ {
				sims = append(sims, page.Sims[i].UniqueName)
			}
		}
	}
	// TODO we could probably get fancy here and send the sims on channels when
	// we get them above, which would give us a head start on fetching the data
	//
	// we could also guard data fetches with a semaphore:
	// github.com/kevinburke/semaphore
	var mu sync.Mutex
	ctx, cancel := context.WithTimeout(context.Background(), 31*time.Second)
	defer cancel()
	group, errctx := errgroup.WithContext(ctx)
	for _, name := range sims {
		name := name
		usage[name] = make([]types.Bits, *duration)
		for i := int(*duration - 1); i >= 0; i-- {
			i := i
			group.Go(func() error {
				start := now.Add(-time.Duration(i+1) * 24 * time.Hour)
				end := now.Add(-time.Duration(i) * 24 * time.Hour)
				page, err := c.Wireless.Sims.GetUsageRecords(errctx, name, url.Values{
					"Start": []string{start.Format(time.RFC3339)},
					"End":   []string{end.Format(time.RFC3339)},
				})
				if err != nil {
					return err
				}
				if len(page.UsageRecords) == 0 {
					return errors.New("no usage records for date range")
				}
				if len(page.UsageRecords) > 1 {
					return errors.New("too many usage records for date range")
				}
				mu.Lock()
				usage[name][i] = page.UsageRecords[0].Data.Total
				mu.Unlock()
				return nil
			})
		}
	}
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
	for _, name := range sims {
		total := types.Bits(0)
		fmt.Printf("%s\n%s\n", name, strings.Repeat("-", len(name)))
		for i := int(*duration - 1); i >= 0; i-- {
			t := now.Add(-time.Duration(i+1) * 24 * time.Hour)
			fmt.Printf("%s: %s\n", t.Format("2006-01-02"), usage[name][i])
			total += usage[name][i]
		}
		fmt.Printf("total (last %d days): %s\n\n", *duration, total)
	}
}
