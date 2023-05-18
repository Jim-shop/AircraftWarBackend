package middleware

import (
	"imshit/aircraftwar/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	// 检查用户登录态
	return func(c *gin.Context) {
		token := common.Token{}
		log.Println(c.ContentType())
		if err := c.ShouldBind(&token); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		//todo
	}
}
