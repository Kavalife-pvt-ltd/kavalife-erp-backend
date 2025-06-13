package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/services"
)

func Routes(r *gin.Engine) {

	apiRoutes := r.Group("/api")                         //default api
	userRoutes := r.Group("/user").Use(AuthMiddleware()) //user api
	productRoutes := r.Group("/product").Use(AuthMiddleware())

	apiRoutes.POST("/login", services.UserLogin)
	apiRoutes.GET("/checkUser", services.CheckUser)
	apiRoutes.POST("/logout", services.Logout)

	userRoutes.GET("/allUsers", services.AllUsers)

	productRoutes.GET("/allProducts", services.AllProducts)
	productRoutes.POST("/insertProduct", services.InsertProduct)
	productRoutes.POST("/updateProduct", services.UpdateProduct)
}
