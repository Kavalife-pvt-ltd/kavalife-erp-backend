package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func CreateGRN(c *gin.Context) {
	var req models.CreateGRNRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid request body"))
		return
	}

	grn, err := handlers.InsertGRNHandler(c, req)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}

	res := models.CreateGRNResponse{
		GRNNumber: grn.GRNNumber,
		CreatedAt: grn.CreatedAt,
	}

	utils.SuccessWithData(c, res)
}

func ViewGRNs(c *gin.Context) {
	grnNo := c.Query("grnNo")

	data, err := handlers.ViewGRNsHandler(c, grnNo)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}

	utils.SuccessWithData(c, data)
}
