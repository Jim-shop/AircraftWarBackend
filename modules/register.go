package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest LoginRequest

func Register(c *gin.Context) {
	// 登记注册信息
	request := RegisterRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//todo
	c.Status(http.StatusOK)
}
