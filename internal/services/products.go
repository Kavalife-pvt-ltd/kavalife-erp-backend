package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllProducts(c *gin.Context) {
	data, err := handlers.AllProductsData(c)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, data)

}

func InsertProduct(c *gin.Context) {
	var req models.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid JSON input"))
		return
	}
	// Get user ID from context
	userID, exists := c.Get("id")
	if !exists {
		utils.BadRequest(c, errors.New("user ID not found in context"))
		return
	}
	uid, ok := userID.(int)
	if !ok {
		utils.BadRequest(c, errors.New("invalid user ID format"))
		return
	}
	req.UserId = uid

	err := handlers.AddProduct(c, req)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}

	utils.SuccessWithData(c, "data inserted")
}
