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

	userRoutes.GET("/allUsers", services.AllUsers)
	userRoutes.POST("/getOneUser", services.GetOneUser)

	productRoutes.GET("/allProducts", services.AllProducts)
	productRoutes.POST("/insertProduct", services.InsertProduct)
	productRoutes.PUT("/updateProduct")
	// api.GET("/authUsersList", services.GetAuthUsers)
}
