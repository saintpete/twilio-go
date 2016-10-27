// +build !go1.7

package twilio

import (
	"net/http"

	"golang.org/x/net/context"
)

func withContext(r *http.Request, ctx context.Context) *http.Request {
	return r
}
