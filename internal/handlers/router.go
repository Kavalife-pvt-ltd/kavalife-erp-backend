package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/services"
)

func Routes(r *gin.Engine) {

	api := r.Group("/api")
	api.POST("/login", services.UserLogin)

	userRoutes := r.Group("/user").Use(AuthMiddleware())
	userRoutes.GET("/allUsers", services.AllUsers)

	// api.GET("/authUsersList", services.GetAuthUsers)
}
