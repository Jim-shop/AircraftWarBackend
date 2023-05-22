package modules

import (
	"imshit/aircraftwar/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	User     string `form:"user" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

// 核验身份并签发token
func Login(c *gin.Context) {
	// 检查提交的信息
	request := LoginRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查用户是否存在
	user, err := models.QueryUser(request.User)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 签发新Token
	token := models.NewToken(user)
	c.String(http.StatusOK, token.Token)
}
