package handlers

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {

	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	api.GET("/allUsers", AllUsers)
}
