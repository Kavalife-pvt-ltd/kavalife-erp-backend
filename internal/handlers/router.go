package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/services"
)

func Routes(r *gin.Engine) {

	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	api.GET("/allUsers", services.AllUsers)
	api.POST("/getOneUser", services.GetOneUser)
	api.GET("/authUsersList", services.GetAuthUsers)
}
