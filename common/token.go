package common

import (
	// "github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Token string `form:"token" json:"token" uri:"token" xml:"token" binding:"required"`
}

func NewToken(account string) string {
	//todo
	return "TOKEN"
}
