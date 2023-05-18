package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	scoreboard := r.Group("scoreboard")
	{
		scoreboard.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		scoreboard.PUT("", func(c *gin.Context) {
			
		})
		scoreboard.POST("/:id", func(c *gin.Context) {
			// id := c.Param("id")
		})
		scoreboard.DELETE("/:id", func(c *gin.Context) {
			// id := c.Param("id")
		})
	}
	login := r.Group("login")
	{
		login.POST("", func(c *gin.Context) {
			
		})
	}
	register := r.Group("register")
	{
		register.POST("", func(c *gin.Context) {
			
		})
	}
	if err := r.Run(":80"); err != nil {
		log.Printf("err: %v\n", err)
		return
	}
}
