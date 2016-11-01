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

type AccessToken struct {
	accountSid string
	apiKey     string
	apiSecret  []byte
	identity   string
	ttl        time.Duration
	NotBefore  time.Time
	grants     map[string]interface{}
	mu         sync.Mutex
}

// To generate an apiKey and apiSecret follow this instructions:
// https://www.twilio.com/docs/api/rest/access-tokens#jwt-format
func New(accountSid, apiKey, apiSecret, identity string, ttl time.Duration) *AccessToken {
	return &AccessToken{
		accountSid: accountSid,
		apiKey:     apiKey,
		apiSecret:  []byte(apiSecret),
		identity:   identity,
		ttl:        ttl,
		mu:         sync.Mutex{},
		grants:     map[string]interface{}{},
	}
}

func (t *AccessToken) AddGrant(grant Grant) {
	t.mu.Lock()
	t.grants["identity"] = t.identity
	t.grants[grant.Key()] = grant.ToPayload()
	t.mu.Unlock()
}

func (t *AccessToken) JWT() (string, error) {
	tNow := time.Now()

	stdClaims := &jwt.StandardClaims{
		Id:        fmt.Sprintf("%s-%d", t.apiKey, tNow.UTC().Unix()),
		ExpiresAt: tNow.Add(t.ttl).UTC().Unix(),
		Issuer:    t.apiKey,
		IssuedAt:  tNow.Unix(),
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
