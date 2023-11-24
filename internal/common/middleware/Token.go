package helpers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims structure to match the expected format of JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"` // Add this field
	jwt.StandardClaims
}

// TokenMiddleware for JWT token validation
func TokenMiddleware(c *gin.Context) {
	// Extract the token from the cookie
	tokenString, err := c.Cookie("token")
	fmt.Printf("Token: %v", tokenString)
	if err != nil {
		// More specific error message
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookie 'token' not found", "details": err.Error()})
		c.Abort()
		return
	}

	// Validate the token
	token, err := ValidateToken(tokenString, "LKJSDFS878dfsdLHLF$lkajd")
	if err != nil {
		// Log the error for debugging
		fmt.Println("Token Validation Error:", err.Error())

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed", "details": err.Error()})
		c.Abort()
		return
	}
	fmt.Println("Token validated successfully") // Debug line
	fmt.Println()

	// Token is valid, extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims extraction failed"})
		c.Abort()
		return
	}

	// Debug: Print out the extracted UserID to see if it's correctly parsed

	c.Set("token", tokenString) // This will make the token available in the context
	c.Set("username", claims.Username)
	c.Set("userID", claims.UserID)
	c.Set("userType", claims.UserType)
}
