package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

func CreateVIR(c *gin.Context) {
	var req models.CreateVIRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid request body"))
		return
	}

	vir, err := handlers.InsertVIRHandler(c, req)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}

	res := models.CreateVIRResponse{
		VIRNumber: vir.VIRNumber,
		CreatedBy: req.CreatedBy.Data.Username,
		CreatedAt: vir.CreatedAt.Format(time.RFC3339),
		Status:    vir.Status,
	}

	utils.SuccessWithData(c, res)
}

func GetAllVIR(c *gin.Context) {
	data, err := handlers.GetAllVIR(c)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}
	utils.SuccessWithData(c, data)
}

func GetVIRByNumber(c *gin.Context) {
	virNumber := c.Param("vir_number")
	if virNumber == "" {
		utils.BadRequest(c, errors.New("missing vir_number"))
		return
	}

	data, err := handlers.GetVIRByNumber(c, virNumber)
	if err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}
	utils.SuccessWithData(c, data)
}

func VerifyVIR(c *gin.Context) {
	virNumber := c.Param("vir_number")
	if virNumber == "" {
		utils.BadRequest(c, errors.New("missing vir_number"))
		return
	}

	var req models.VerifyVIRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, errors.New("invalid request body"))
		return
	}

	checkedAtParsed, err := time.Parse(time.RFC3339, req.CheckedAt)
	if err != nil {
		utils.BadRequest(c, errors.New("invalid checkedAt format"))
		return
	}

	if err := handlers.UpdateVIRVerification(c, virNumber, req); err != nil {
		utils.StatusInternalServerError(c, err)
		return
	}

	res := models.VerifyVIRResponse{
		VIRNumber: virNumber,
		CheckedBy: req.CheckedBy.Data.Username,
		CheckedAt: checkedAtParsed.Format(time.RFC3339),
		Status:    "completed",
	}

	utils.SuccessWithData(c, res)
}
