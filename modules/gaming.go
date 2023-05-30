package modules

import (
	"imshit/aircraftwar/models"
	"imshit/aircraftwar/socket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

func Socket(c *gin.Context) {
	// 获取 Token
	token := &models.Token{}
	if err := c.ShouldBind(token); err != nil {
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
	// 获取用户
	user, err := models.GetUser(userId)
	if err != nil {
		log.Printf("Get user error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket estabulish error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 进入配对阶段
	socket.Pairing(ws, user)
}
