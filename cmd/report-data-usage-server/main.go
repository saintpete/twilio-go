package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	types "github.com/kevinburke/go-types"
	"github.com/kevinburke/handlers"
	"github.com/kevinburke/rest"
	twilio "github.com/kevinburke/twilio-go"
	"github.com/kevinburke/twilio-go/datausage"
	"golang.org/x/sync/errgroup"
)

const Version = "0.1"

var duration = flag.Uint("days", 7, "Number of days to get usage for")
var location = flag.String("location", "", "Timezone to use (defaults to system location/TZ env var)")
var sim = flag.String("sim", "", "Load data about a specific Sim/unique name only")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Start a server that prints information about how much data you used on each day.
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
	iter := c.Wireless.Sims.GetPageIterator(nil)
	sims := make([]string, 0)
	if *sim != "" {
		sims = []string{*sim}
	} else {
		// TODO, maybe need to look this up on a loop every hour, or something.
		// for now just running this as a demo.
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		days := *duration
		if d := r.URL.Query().Get("days"); d != "" {
			daysInt, err := strconv.ParseUint(d, 10, 64)
			if err == nil {
				days = uint(daysInt)
			}
		}
		end := time.Now().Add(24 * time.Hour).In(loc)
		end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
		start := end.Add(-24 * time.Hour * time.Duration(days+1))
		var mu sync.Mutex
		ctx, cancel := context.WithTimeout(r.Context(), 31*time.Second)
		defer cancel()
		group, errctx := errgroup.WithContext(ctx)
		usage := make(map[string][]types.Bits)
		for _, name := range sims {
			name := name
			group.Go(func() error {
				u, err := datausage.GetUsage(errctx, c, name, start, end, 24*time.Hour)
				if err != nil {
					return err
				}
				mu.Lock()
				usage[name] = u
				mu.Unlock()
				return nil
			})
		}
		if err := group.Wait(); err != nil {
			rest.ServerError(w, r, err)
		}
		buf := new(bytes.Buffer)
		for _, name := range sims {
			total := types.Bits(0)
			fmt.Fprintf(buf, "%s\n%s\n", name, strings.Repeat("-", len(name)))
			for i := 0; i < len(usage[name]); i++ {
				t := start.Add(time.Duration(i) * 24 * time.Hour)
				fmt.Fprintf(buf, "%s: %s\n", t.Format("2006-01-02"), usage[name][i])
				total += usage[name][i]
			}
			fmt.Fprintf(buf, "total (last %d days): %s\n\n", days, total)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(buf.Bytes())
	})
	addr := ":4425"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		handlers.Logger.Error("Error listening", "addr", addr, "err", err)
		os.Exit(2)
	}
	mux := handlers.UUID(http.DefaultServeMux) // add UUID header
	mux = handlers.Server(mux,
		"report-data-usage-server/"+Version) // add Server header
	mux = handlers.Log(mux)      // log requests/responses
	mux = handlers.Duration(mux) // add Duration header
	handlers.Logger.Info("Started server", "protocol", "http", "port", addr)
	http.Serve(ln, mux)
}
