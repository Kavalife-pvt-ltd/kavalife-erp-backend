package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func AllVendors(c *gin.Context) {
	data, err := handlers.AllVendorsData(c)
	if err != nil {
		utils.SuccessWithError(c, err)
		return
	}
	utils.SuccessWithData(c, data)
}

func InsertVendors(c *gin.Context) {
	var req models.Vendors
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

	err := handlers.AddVendor(c, req, uid)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}

	utils.SuccessWithData(c, "data inserted")
}

// func UpdateVendor(c *gin.Context) {

// 	req := models.VendorUpdate
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		utils.BadRequest(c, errors.New("invalid input data"))
// 		return
// 	}
// 	// Get user ID from context
// 	userID, exists := c.Get("id")
// 	if !exists {
// 		utils.BadRequest(c, errors.New("user ID not found in context"))
// 		return
// 	}
// 	uid, ok := userID.(int)
// 	if !ok {
// 		utils.BadRequest(c, errors.New("invalid user ID format"))
// 		return
// 	}

// 	// Call DB function to update
// 	err := handlers.UpdateVendor(c, req.ID, req.Quantity, uid)
// 	if err != nil {
// 		utils.StatusInternalServerError(c, err)
// 		return
// 	}

// 	utils.SuccessWithData(c, "product updated successfully")
// }
