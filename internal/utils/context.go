package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paaart/kavalife-erp-backend/internal/models"
)

func GetUserID(c *gin.Context) (int64, bool) {
	// Try canonical key first
	if val, exists := c.Get("user_id"); exists {
		return extractID(val)
	}

	// Fallback to "id" (your middleware)
	if val, exists := c.Get("id"); exists {
		return extractID(val)
	}

	return 0, false
}

// internal helper to convert any type into int64
func extractID(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case float64:
		return int64(v), true
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n, true
		}
		return 0, false
	case models.User:
		return int64(v.ID), true
	case *models.User:
		if v != nil {
			return int64(v.ID), true
		}
		return 0, false
	default:
		fmt.Printf("Unexpected user ID type: %T\n", val)
		return 0, false
	}
}

func GetUserRole(c *gin.Context) (string, bool) {
	if val, exists := c.Get("role"); exists {
		if s, ok := val.(string); ok {
			return s, true
		}
	}
	return "", false
}
