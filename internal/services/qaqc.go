package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func CreateQAQC(c *gin.Context) {
	var req models.CreateQAQCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := handlers.InsertQAQCHandler(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert QAQC entry", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.CreateQAQCResponse{
		Message: "QA/QC entry created successfully",
	})
}

func ViewQAQC(c *gin.Context) {
	processType := c.Query("processType")
	processRef := c.Query("processRef")

	if processType == "" || processRef == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing processType or processRef"})
		return
	}

	entry, err := handlers.GetQAQCByRef(processType, processRef)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch QAQC", "details": err.Error()})
		return
	}
	if entry == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No QAQC entry found for this reference"})
		return
	}

	c.JSON(http.StatusOK, entry)
}
