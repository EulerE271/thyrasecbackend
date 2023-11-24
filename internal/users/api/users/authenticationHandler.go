package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"thyra/internal/common/db"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("LKJSDFS878dfsdLHLF$lkajd")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID   string `json:"user_id"` // Change to string for UUID
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}

func LoginHandler(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	dbConn := db.GetDB() // Use GetDB() to obtain the database connection

	userTypes := map[string]string{
		"admins":            "admin",
		"partners_advisors": "partner_advisor",
		"customers":         "customer",
	}

	var userType string
	var storedPassword string
	var err error

	var userID string

	for table, uType := range userTypes {
		query := fmt.Sprintf("SELECT id, password_hash FROM %s WHERE username = $1", table)
		err = dbConn.QueryRow(query, loginRequest.Username).Scan(&userID, &storedPassword)
		if err == nil {
			userType = uType
			log.Printf("User found in %s table with ID %s", table, userID)
			break
		} else if err != sql.ErrNoRows {
			log.Printf("Error querying %s table: %v", table, err)
			break // Exit the loop on error other than ErrNoRows
		}
	}

	if err == sql.ErrNoRows || userType == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Assuming you have a function generateJWTToken that creates the JWT token
	token, err := generateJWTToken(userID, loginRequest.Username, userType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	fmt.Println("This is the token, right before it is set: %v", token)
	c.SetCookie("token", token, 86400, "/", "http://dev.thyrasolutions.se", true, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateJWTToken(userID string, username, userType string) (string, error) { // Accept userID as a parameter
	claims := Claims{
		UserID:   userID, // <-- Set the UserID
		Username: username,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
