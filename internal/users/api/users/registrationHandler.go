package handlers

import (
	"net/http"
	"thyra/internal/common/db"
	customerno "thyra/internal/users/services/customernumber"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Common attributes for all user types
type BaseRegistrationRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	CustomerNumber string `json:"customer_number"`
}

// For Admin
type AdminRegistrationRequest struct {
	BaseRegistrationRequest
	// Add any additional fields unique to admins if needed
}

// For Partners/Advisors
type PartnerAdvisorRegistrationRequest struct {
	BaseRegistrationRequest
	FullName    string `json:"full_name"`
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
}

// For Customers
type CustomerRegistrationRequest struct {
	BaseRegistrationRequest
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// Adjusting the registerUser function to return an error
func insertIntoDb(query string, args ...interface{}) error {
	dbConn := db.GetDB() // Correctly use GetDB()
	_, err := dbConn.Exec(query, args...)
	return err
}
func RegisterAdminHandler(c *gin.Context) {
	var registrationRequest AdminRegistrationRequest
	if err := c.BindJSON(&registrationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	customerNumber := customerno.GenerateCustomerNumber()

	query := "INSERT INTO admins (username, password_hash, email, customer_number) VALUES ($1, $2, $3, $4)"
	if err := insertIntoDb(query, registrationRequest.Username, string(hashedPassword), registrationRequest.Email, customerNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting admin data into the database", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful for admin"})
}

func RegisterPartnerAdvisorHandler(c *gin.Context) {
	var registrationRequest PartnerAdvisorRegistrationRequest
	if err := c.BindJSON(&registrationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	customerNumber := customerno.GenerateCustomerNumber()

	query := "INSERT INTO partners_advisors (username, password_hash, email, full_name, company_name, phone_number, customer_number) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if err := insertIntoDb(query, registrationRequest.Username, string(hashedPassword), registrationRequest.Email, registrationRequest.FullName, registrationRequest.CompanyName, registrationRequest.PhoneNumber, customerNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting partner/advisor data into the database", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful for partner/advisor"})
}

func RegisterCustomerHandler(c *gin.Context) {
	var registrationRequest CustomerRegistrationRequest
	if err := c.BindJSON(&registrationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	customerNumber := customerno.GenerateCustomerNumber()

	query := "INSERT INTO customers (username, password_hash, email, full_name, address, phone_number, customer_number) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if err := insertIntoDb(query, registrationRequest.Username, string(hashedPassword), registrationRequest.Email, registrationRequest.FullName, registrationRequest.Address, registrationRequest.PhoneNumber, customerNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting customer data into the database", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful for customer"})
}

func CookieTestHandler(c *gin.Context) {
	// Set a test cookie
	c.SetCookie("testcookiereghandler", "testvaluereghandler", 86400, "/", "dev.thyrasolutions.se", true, false)

	// Send a simple response
	c.JSON(http.StatusOK, gin.H{"message": "Test cookie set"})
}
