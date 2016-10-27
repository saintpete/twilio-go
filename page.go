package twilio

import (
	"errors"
	"net/url"

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
