package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetScoreboard(c *gin.Context) {
	// 获取总排行榜
	//todo
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func AddScoreboard(c *gin.Context) {
	// 新增排行榜项目
	//todo
}

func DeleteScoreboard(c *gin.Context) {
	// 删除排行榜项目
	//todo
	// id := c.Param("id")
}
