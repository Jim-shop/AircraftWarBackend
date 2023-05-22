package modules

import (
	"imshit/aircraftwar/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest LoginRequest

func Register(c *gin.Context) {
	// 获取注册信息
	request := RegisterRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 核验是否符合要求
	if size := len(request.User); size <= 0 && size >= 16 {
		c.Status(http.StatusBadRequest)
		return
	}
	if size := len(request.Password); size <= 0 && size >= 16 {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查用户是否存在
	if _, err := models.QueryUser(request.User); err == nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 添加新用户
	if err := models.CreateUser(&models.User{
		Name:     request.User,
		Password: request.Password,
	}); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 返回成功标志
	c.Status(http.StatusOK)
}
