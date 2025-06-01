package handlers

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
