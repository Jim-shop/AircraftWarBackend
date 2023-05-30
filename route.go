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
	gaming := r.Group("gaming", middleware.AuthMiddleWare())
	{
		gaming.GET("", modules.Socket)
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
