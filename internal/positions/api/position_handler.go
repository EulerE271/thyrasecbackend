package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPositionsById(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to create an instrument"})
		return
	}

}
