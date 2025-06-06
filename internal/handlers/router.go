package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/services"
)

func Routes(r *gin.Engine) {

	api := r.Group("/api")
	userRoutes := r.Group("/user") //user api
	userRoutes.Use(AuthMiddleware())
	userRoutes.GET("/allUsers", services.AllUsers)
	userRoutes.POST("/getOneUser", services.GetOneUser)
	api.POST("/login", services.UserLogin)
	// api.GET("/authUsersList", services.GetAuthUsers)
}
