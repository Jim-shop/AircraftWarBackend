package modules

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"imshit/aircraftwar/daemon"
	"github.com/spf13/viper"
	"imshit/aircraftwar/models"
)

type PairingRequest struct {
	models.Token
	Mode  string    `form:"mode" json:"mode" uri:"mode" xml:"mode" binding:"required"`
}

func Pairing(c *gin.Context) {
	// 检验请求
	request := &PairingRequest{}
	if err := c.ShouldBind(request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Mode != "easy" && request.Mode != "medium" && request.Mode != "hard" {
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户ID
	userId, err := models.GetUserIDByToken(&request.Token)
	if err != nil {
		log.Printf("Get user id error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户
	user, err := models.GetUser(userId)
	if err != nil {
		log.Printf("Get user error: %v\n", err)
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
	daemon.GetPairingDaemon().Bind(ws, user, request.Mode)
}
