package handlers

import (
	"database/sql"
	"net/http"
	"thyra/internal/common/db"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	UUID           string `json:"uuid"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	CustomerNumber string `json:"customer_number"`
	// Add any other fields you want to be returned in the API response
}

func GetAllUsersHandler(c *gin.Context) {
	role := c.DefaultQuery("role", "customer") // default to "customer" if no role provided
	var query string

	dbConn := db.GetDB() // Use GetDB() to obtain the database connection

	switch role {
	case "admin":
		query = "SELECT id, username, email, customer_number FROM admins"
	case "advisor":
		query = "SELECT id, username, email, customer_number FROM partners_advisors"
	case "customer":
		query = "SELECT id, username, email, customer_number FROM customers"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	rows, err := dbConn.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users from the database", "detail": err.Error()})
		return
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var user UserResponse
		err = rows.Scan(&user.UUID, &user.Username, &user.Email, &user.CustomerNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row into struct", "detail": err.Error()})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing rows", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserNameByUuid(c *gin.Context) {
	// Extract the authenticated user's ID and role from context.
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

	// Ensure the authenticated user is an admin.
	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to fetch a username"})
		return
	}

	// Extract the UUID from the query parameter.
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID parameter is required"})
		return
	}

	// Get the db connection
	dbConn := db.GetDB() // Correctly use GetDB()

	// Query the database to get the username associated with the UUID.
	var username string
	query := `SELECT full_name FROM customers WHERE id = $1`
	err := dbConn.QueryRow(query, uuid).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch username", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": username})
}
