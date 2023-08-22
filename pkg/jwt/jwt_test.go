package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestJWT(t *testing.T) {
	userJwt := NewJWT([]byte("signingKey"))

	token, err := userJwt.CreateToken(CustomClaims{
		1,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "dousheng",
		},
	})
	fmt.Println(token)
	if err != nil {
		t.Fatalf("create token error %v", err)
	}

	claims, err := userJwt.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MSwiZXhwIjoxNjkyNzczOTEzLCJpc3MiOiJkb3VzaGVuZyJ9.NyBKyKiRM7uopjtJ-d0PUALLR6V71dGS5-oyUqfCO6k")
	fmt.Println(claims.Id)

	//if err != nil {
	//	t.Fatalf("token verified error %v", err)
	//}
	//
	//otherJwt := NewJWT([]byte{0x12, 0x32, 0x4a, 0x53, 0x59, 0x45})
	//_, err = otherJwt.ParseToken(token)
	//
	//if err != nil {
	//	t.Fatalf("token verified error %v", err)
	//}
	//
	//time.Sleep(time.Second * 7)
	//
	//_, err = userJwt.ParseToken(token)
	//
	//if err == nil {
	//	t.Fatalf("token expired but not got error")
	//}

}
