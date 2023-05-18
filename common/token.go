package common

type Token struct {
	Token string `form:"token" json:"token" uri:"token" xml:"token" binding:"required"`
}