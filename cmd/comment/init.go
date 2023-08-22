package main

import (
	"github.com/alph00/tiktok-tiny/pkg/jwt"
)

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}
