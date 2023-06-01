/**
 * Copyright (c) [2023] [Jim-shop]
 * [AircraftWarBackend] is licensed under Mulan PubL v2.
 * You can use this software according to the terms and conditions of the Mulan PubL v2.
 * You may obtain a copy of Mulan PubL v2 at:
 *          http://license.coscl.org.cn/MulanPubL-2.0
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PubL v2 for more details.
 */

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

type PairingRequest struct {
	models.Token
	Mode string `form:"mode" json:"mode" uri:"mode" xml:"mode" binding:"required"`
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
