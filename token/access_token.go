// Package token generates valid tokens for Twilio Client SDKs.

// Create them on your server and pass them to a client to help verify a
// client's identity, and to grant access to features in client API's.
//
// For more information, please see the Twilio docs:
// https://www.twilio.com/docs/api/rest/access-tokens
package token

import (
	"fmt"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const jwtContentType = "twilio-fpa;v=1"

type myCustomClaims struct {
	Grants map[string]interface{} `json:"grants"`
	*jwt.StandardClaims
}

// AccessToken holds properties that can generate a signature to talk to the
// Messaging/Voice/Video API's.
type AccessToken struct {
	NotBefore  time.Time
	accountSid string
	apiKey     string
	apiSecret  []byte
	identity   string
	ttl        time.Duration
	grants     map[string]interface{}
	mu         sync.Mutex
}

// New generates a new AccessToken with the specified properties. identity is
// a unique ID for a particular user.
//
// To generate an apiKey and apiSecret follow this instructions:
// https://www.twilio.com/docs/api/rest/access-tokens#jwt-format
func New(accountSid, apiKey, apiSecret, identity string, ttl time.Duration) *AccessToken {
	return &AccessToken{
		accountSid: accountSid,
		apiKey:     apiKey,
		apiSecret:  []byte(apiSecret),
		identity:   identity,
		ttl:        ttl,
		grants:     make(map[string]interface{}),
	}
}

func (t *AccessToken) AddGrant(grant Grant) {
	t.mu.Lock()
	t.grants["identity"] = t.identity
	t.grants[grant.Key()] = grant.ToPayload()
	t.mu.Unlock()
}

func (t *AccessToken) JWT() (string, error) {
	now := time.Now().UTC()

	stdClaims := &jwt.StandardClaims{
		Id:        fmt.Sprintf("%s-%d", t.apiKey, now.Unix()),
		ExpiresAt: now.Add(t.ttl).Unix(),
		Issuer:    t.apiKey,
		IssuedAt:  now.Unix(),
		Subject:   t.accountSid,
	}
	if !t.NotBefore.IsZero() {
		stdClaims.NotBefore = t.NotBefore.UTC().Unix()
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	claims := myCustomClaims{
		Grants:         t.grants,
		StandardClaims: stdClaims,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken.Header["cty"] = jwtContentType
	return jwtToken.SignedString(t.apiSecret)
}
