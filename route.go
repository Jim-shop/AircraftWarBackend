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

package main

import (
	"imshit/aircraftwar/middleware"
	"imshit/aircraftwar/modules"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupServer() (r *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()
	login := r.Group("login")
	{
		login.POST("", modules.Login)
	}
	register := r.Group("register")
	{
		register.POST("", modules.Register)
	}
	game := r.Group("game", middleware.AuthMiddleWare())
	{
		game.GET("pairing", modules.Pairing)
		game.GET("fighting/:id", modules.Fighting)
	}
	scoreboard := r.Group("scoreboard", middleware.AuthMiddleWare())
	{
		scoreboard.GET("", modules.GetScoreboard)
		scoreboard.PUT("", modules.AddScoreboard)
		scoreboard.DELETE("/:id", modules.DeleteScoreboard)
	}
	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})
	return r
}
