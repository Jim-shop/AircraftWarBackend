package modules

import (
	"imshit/aircraftwar/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	if size := len(request.User); size < 2 && size > 16 {
		c.Status(http.StatusBadRequest)
		return
	}
	if size := len(request.Password); size != 64 {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查用户是否存在
	if _, err := models.QueryUser(request.User); err == nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 加密存储密码
	crypt, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Bcrypt error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 添加新用户
	if err := models.CreateUser(&models.User{
		Name:     request.User,
		Password: crypt,
	}); err != nil {
		log.Printf("Create user error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户
	user, err := models.QueryUser(request.User)
	if err != nil {
		log.Printf("Query user error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 签发Token
	token, err := models.NewToken(user, c)
	if err != nil {
		log.Printf("Token generate error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.String(http.StatusOK, token.Token)
}
