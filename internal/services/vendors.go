package services

import (
	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
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
