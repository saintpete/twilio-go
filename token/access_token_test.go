package token

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	ACC_SID    = "123456"
	API_KEY    = "abcdef"
	API_SECRET = "asdfghjklqwertyuiopzxcvbnm"
	IDENTITY   = "johnsmith"
	APP_SID    = "asdfghjkl"
)

func TestJWT(t *testing.T) {
	t.Parallel()

	accTkn := New(ACC_SID, API_KEY, API_SECRET, IDENTITY, time.Hour)
	accTkn.NotBefore = time.Now()
	convGrant := NewConversationsGrant(APP_SID)

	accTkn.AddGrant(convGrant)
	jwtString, err := accTkn.JWT()

	if err != nil {
		t.Error("Unexpected error when generating the token", err)
	}
	if jwtString == "" {
		t.Error("token returned is empty")
	}

	token, err := jwt.ParseWithClaims(jwtString, &myCustomClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(API_SECRET), nil
	})
	if err != nil {
		t.Error("Unexpected error when generating the token", err)
	}

	claims := token.Claims.(*myCustomClaims)

	if &claims.StandardClaims == nil {
		t.Error("Claim doesn't conaint a standard claims struct")
	}

	if claims.StandardClaims.ExpiresAt == 0 {
		t.Error("ExpiredAt is not set")
	}

	if claims.StandardClaims.Id == "" {
		t.Error("ID is not set")
	}

	if claims.StandardClaims.IssuedAt == 0 {
		t.Error("IssuedAt is not set")
	}

	if claims.StandardClaims.NotBefore == 0 {
		t.Error("NotBefore is not set")
	}

	if claims.StandardClaims.Issuer != API_KEY {
		t.Errorf("Issuer expected to be: %s, got %s\n", API_KEY, claims.StandardClaims.Issuer)
	}

	if claims.StandardClaims.Subject != ACC_SID {
		t.Errorf("Subject expected to be: %s, got %s\n", ACC_SID, claims.StandardClaims.Subject)
	}

	if claims.Grants == nil {
		t.Error("Expected Grants to exist")
	}

	if claims.Grants["identity"] != IDENTITY {
		t.Errorf("Grants identity expected to be %s, got %s\n", IDENTITY, claims.Grants["identity"])
	}
}
