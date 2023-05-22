package models

import (
	// "github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Token string `form:"token" json:"token" uri:"token" xml:"token" binding:"required"`
}

func NewToken(user *User) *Token {
	//todo
	return &Token{"TOKEN"}
}

func ValidateToken(token *Token) bool {
	// todo
	return true
}
