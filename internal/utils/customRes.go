package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func SuccessWithError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"error": err.Error(),
	})
}
func BadRequest(c *gin.Context, err error) { //400
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}
func StatusUnauthorized(c *gin.Context, err error) { //401
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": err.Error(),
	})
}
func StatusInternalServerError(c *gin.Context, err error) { //500
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
