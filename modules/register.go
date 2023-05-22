package modules

import (
	"imshit/aircraftwar/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest LoginRequest

func Register(c *gin.Context) {
	// 登记注册信息
	request := RegisterRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查用户是否存在
	if _, err := models.QueryUser(request.User); err == nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 添加新用户
	if err := models.AddUser(&models.User{
		Name: request.User,
		Password: request.Password,
	}); err != nil {
		c.Status(http.StatusBadRequest)
		return 
	}
	// 返回成功标志
	c.Status(http.StatusOK)
}
