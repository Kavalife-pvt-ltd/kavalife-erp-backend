package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	"github.com/paaart/kavalife-erp-backend/internal/models"
	"github.com/paaart/kavalife-erp-backend/internal/utils"
)

// POST /sales-po/create
// Create a new sales PO (quote request) from the Information page
func CreateSalesPO(c *gin.Context) {
	var req models.CreateSalesPORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"details": err.Error(),
		})
		return
	}

	// Always trust authenticated user as sales rep
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	req.SalesRepID = &userID

	po, err := handlers.CreateSalesPO(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_create_sales_po",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    po,
	})
}

// GET /sales-po/:id
// Fetch a single PO by id
func GetSalesPO(c *gin.Context) {
	idStr := c.Param("id")
	poID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || poID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_po_id",
		})
		return
	}

	po, err := handlers.GetSalesPOByID(c, poID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "sales_po_not_found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    po,
	})
}

// GET /sales-po/view
// List POs with optional filters:
//
//	?status=quote_requested
//	?salesRepId=123
//	?productId=45
func ListSalesPO(c *gin.Context) {
	status := c.Query("status")
	salesRepIDStr := c.Query("salesRepId")
	productIDStr := c.Query("productId")

	var filter handlers.SalesPOFilter

	if status != "" {
		filter.Status = &status
	}

	if salesRepIDStr != "" {
		if id, err := strconv.ParseInt(salesRepIDStr, 10, 64); err == nil {
			filter.SalesRepID = &id
		}
	}

	if productIDStr != "" {
		if id, err := strconv.ParseInt(productIDStr, 10, 64); err == nil {
			filter.ProductID = &id
		}
	}

	pos, err := handlers.ListSalesPO(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_list_sales_po",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pos,
	})
}

// PATCH /sales-po/:id/status
// Update status (admin/client transitions) — approve, reject, route, etc.
func UpdateSalesPOStatus(c *gin.Context) {
	idStr := c.Param("id")
	poID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || poID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_po_id",
		})
		return
	}

	var req models.UpdateSalesPOStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"details": err.Error(),
		})
		return
	}
	req.POID = poID

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	po, err := handlers.UpdateSalesPOStatus(c, req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed_to_update_status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    po,
	})
}

// GET /sales-po/:id/status-log
// Fetch full status history for a PO
func GetSalesPOStatusLog(c *gin.Context) {
	idStr := c.Param("id")
	poID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || poID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_po_id",
		})
		return
	}

	logs, err := handlers.GetSalesPOStatusLog(c, poID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_fetch_status_log",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
	})
}
