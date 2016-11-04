package twilio

import (
	"errors"
	"net/url"
	"strings"
	"time"

	types "github.com/kevinburke/go-types"
	"golang.org/x/net/context"
)

type Page struct {
	FirstPageURI    string           `json:"first_page_uri"`
	Start           uint             `json:"start"`
	End             uint             `json:"end"`
	NumPages        uint             `json:"num_pages"`
	Total           uint             `json:"total"`
	NextPageURI     types.NullString `json:"next_page_uri"`
	PreviousPageURI types.NullString `json:"previous_page_uri"`
	PageSize        uint             `json:"page_size"`
}

type Meta struct {
	FirstPageURL    string           `json:"first_page_url"`
	NextPageURL     types.NullString `json:"next_page_url"`
	PreviousPageURL types.NullString `json:"previous_page_url"`
	Key             string           `json:"key"`
	Page            uint             `json:"page"`
	PageSize        uint             `json:"page_size"`
}

// NoMoreResults is returned if you reach the end of the result set while
// paging through resources.
var NoMoreResults = errors.New("twilio: No more results")

type PageIterator struct {
	client      *Client
	nextPageURI types.NullString
	data        url.Values
	count       uint
	pathPart    string
}

func (p *PageIterator) SetNextPageURI(npuri types.NullString) {
	if npuri.Valid == false {
		p.nextPageURI = npuri
		return
	}
	if strings.HasPrefix(npuri.String, p.client.Base) {
		npuri.String = npuri.String[len(p.client.Base):]
	}
	p.nextPageURI = npuri
}

// Next asks for the next page of resources and decodes the results into v.
func (p *PageIterator) Next(ctx context.Context, v interface{}) error {
	var err error
	if p.count == 0 {
		err = p.client.ListResource(ctx, p.pathPart, p.data, v)
	} else if p.nextPageURI.Valid == false {
		return NoMoreResults
	} else {
		err = p.client.GetNextPage(ctx, p.nextPageURI.String, v)
	}
	if err != nil {
		return err
	}
	p.count++
	return nil
}

func NewPageIterator(client *Client, data url.Values, pathPart string) *PageIterator {
	return &PageIterator{
		data:        data,
		client:      client,
		count:       0,
		nextPageURI: types.NullString{},
		pathPart:    pathPart,
	}
}

// containsResultsInRange returns true if any results are in the range
// [start, end).
func containsResultsInRange(start time.Time, end time.Time, results []time.Time) bool {
	for _, result := range results {
		if (result.Equal(start) || result.After(start)) && result.Before(end) {
			return true
		}
	}
	return false
}

// shouldContinuePaging returns true if fetching more results (that have
// earlier timestamps than the provided results) could possibly return results
// in the range. shouldContinuePaging assumes results is sorted so the first
// result in the slice has the latest timestamp, and the last result in the
// slice has the earliest timestamp. shouldContinuePaging panics if results is
// empty.
func shouldContinuePaging(start time.Time, results []time.Time) bool {
	// the last result in results is the earliest. if the earliest result is
	// before the start, fetching more resources may return more results.
	if len(results) == 0 {
		panic("zero length result set")
	}
	last := results[len(results)-1]
	return last.After(start)
}

// indexesOutsideRange returns the indexes of times in results that are outside
// of [start, end). indexesOutsideRange panics if start is later than end.
func indexesOutsideRange(start time.Time, end time.Time, results []time.Time) []int {
	if start.After(end) {
		panic("start date is after end date")
	}
	indexes := make([]int, 0, len(results))
	for i, result := range results {
		if result.Equal(end) || result.After(end) {
			indexes = append(indexes, i)
		}
		if result.Before(start) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}
