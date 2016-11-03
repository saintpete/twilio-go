package twilioclient

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestCapability(t *testing.T) {
	t.Parallel()
	cap := NewCapability("AC123", "123")
	cap.AllowClientIncoming("client-name")
	tok, err := cap.GenerateToken(time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tok)
	cc := new(customClaim)
	_, err = jwt.ParseWithClaims(tok, cc, func(tkn *jwt.Token) (interface{}, error) {
		return []byte("123"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if cc.StandardClaims.Issuer != "AC123" {
		t.Errorf("bad Issuer")
	}
	if cc.Scope != "scope:client:incoming?clientName=client-name" {
		t.Errorf("bad Scope")
	}
}
