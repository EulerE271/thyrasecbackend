package utils

import (
	"github.com/gin-gonic/gin"
)

func GetAuthenticatedUser(c *gin.Context) (string, string, bool) {
	userID, exists := c.Get("userID")
	if !exists || userID == nil {
		return "", "", false
	}

	userRole, exists := c.Get("userType")
	if !exists || userRole == nil {
		return "", "", false
	}

	return userID.(string), userRole.(string), true
}
