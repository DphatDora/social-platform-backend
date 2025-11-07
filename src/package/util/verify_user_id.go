package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Helper function to extract userID from context
func GetUserIDFromContext(c *gin.Context) (uint64, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("userID not found in context")
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		return 0, fmt.Errorf("invalid userID type")
	}

	return userID, nil
}

func GetOptionalUserIDFromContext(c *gin.Context) *uint64 {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return nil
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		return nil
	}

	return &userID
}
