package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	User     string `form:"user" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

func Login(c *gin.Context) {
	// 核验身份并签发token
	request := LoginRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//todo
	c.String(http.StatusOK, "TOKEN")

}
