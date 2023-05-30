package modules

import (
	"imshit/aircraftwar/daemon"
	"imshit/aircraftwar/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func Fighting(c *gin.Context) {
	// 检验请求
	token := &models.Token{}
	if err := c.ShouldBind(token); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	roomId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if roomId < 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户ID
	userId, err := models.GetUserIDByToken(token)
	if err != nil {
		log.Printf("Get user id error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取WebSocket连接
	upgrader := websocket.Upgrader{
		ReadBufferSize:  viper.GetInt("socket.readBufferSize"),
		WriteBufferSize: viper.GetInt("socket.writeBufferSize"),
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket estabulish error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	daemon.GetFightingDaemon().Connect(ws, userId, int(roomId))
}
