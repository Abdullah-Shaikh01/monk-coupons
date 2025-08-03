package utils

import (
	"github.com/gin-gonic/gin"
)

// Success sends a standardized success response with optional data
func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}

// SuccessMessage sends a success response without a data field
func SuccessMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"message": message,
	})
}

// Error sends a standardized error response
func Error(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"message": message,
		"error":   err.Error(),
	})
}

// Error sends a standardized error response without err
func ErrorWithoutErr(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  statusCode,
		"message": message,
	})
}
