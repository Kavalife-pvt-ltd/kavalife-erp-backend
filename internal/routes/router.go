package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/services"
)

func Routes(r *gin.Engine) {

	apiRoutes := r.Group("/api")                               //default api
	userRoutes := r.Group("/user").Use(AuthMiddleware())       //user api
	productRoutes := r.Group("/product").Use(AuthMiddleware()) //products api
	vendorRoutes := r.Group("/vendor").Use(AuthMiddleware())   //vendor api
	virRoutes := r.Group("/vir").Use(AuthMiddleware())
	grnRoutes := r.Group("/grn").Use(AuthMiddleware())
	qaqc := r.Group("/qaqc").Use(AuthMiddleware())
	salesPORoutes := r.Group("/sales-po").Use(AuthMiddleware()) // 🔹 new

	apiRoutes.POST("/login", services.UserLogin)
	apiRoutes.GET("/checkUser", services.CheckUser)
	// apiRoutes.POST("/logout", services.Logout)

	userRoutes.GET("/allUsers", services.AllUsers)
	userRoutes.GET("/logout", services.Logout)

	productRoutes.GET("/allProducts", services.AllProducts)
	productRoutes.POST("/insertProduct", services.InsertProduct)
	productRoutes.POST("/updateProduct", services.UpdateProduct)

	vendorRoutes.GET("allVendors", services.AllVendors)
	vendorRoutes.POST("/insertVendor", services.InsertVendors)
	// productRoutes.POST("/updateVendor", services.UpdateVendor)

	virRoutes.POST("/create", services.CreateVIR)
	virRoutes.GET("/all", services.GetAllVIR)
	virRoutes.GET("/:vir_number", services.GetVIRByNumber)
	virRoutes.PATCH("/verify/:vir_number", services.VerifyVIR)

	grnRoutes.POST("/create", services.CreateGRN)
	grnRoutes.GET("/view", services.ViewGRNs)

	qaqc.POST("/create", services.CreateQAQC)
	qaqc.GET("/view", services.ViewQAQC)

	salesPORoutes.POST("/create", services.CreateSalesPO)
	salesPORoutes.GET("/view", services.ListSalesPO)
	salesPORoutes.GET("/:id", services.GetSalesPO)
	salesPORoutes.PATCH("/:id/status", services.UpdateSalesPOStatus)
	salesPORoutes.GET("/:id/status-log", services.GetSalesPOStatusLog)
}
