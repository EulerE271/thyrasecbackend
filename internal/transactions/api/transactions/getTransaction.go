package handlers

import (
	"net/http"
	"thyra/internal/common/db"
	"thyra/internal/transactions/models"
	"thyra/internal/transactions/utils"

	"github.com/gin-gonic/gin"
)

//Function for fetching all transactions for a specific user
func GetTransactionByUserHandler(c *gin.Context) {
	// Extracting user parameters
	targetUserID := c.Param("userId")
	authUserID, authUserRole, isAuthenticated := utils.GetAuthenticatedUser(c)

	// Authentication and Authorization checks
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or role not found"})
		return
	}

	if authUserRole != "admin" && authUserID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to fetch another user's transactions"})
		return
	}

	// Database connection
	db := db.GetConnection(c)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection issue"})
		return
	}

	// Fetching transactions for the user directly within the handler
	query := `
    SELECT
        t.*,
        tt.transaction_type_name,
        ts.status_label,
        a.account_number
    FROM 
        thyrasec.transactions t
    LEFT JOIN 
        thyrasec.transactions_types tt ON t.type = tt.type_id
    LEFT JOIN 
        thyrasec.transaction_status ts ON t.status_transaction = ts.status_id
    LEFT JOIN
        thyrasec.accounts a ON t.account_owner_id = a.id
    WHERE 
        t.account_owner_id = $1;
    `
	var transactions []models.TransactionDisplay
	if err := db.Select(&transactions, query, targetUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions", "details": err.Error()})
		return
	}

	// Handling empty transactions
	if len(transactions) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No transactions found"})
		return
	}

	// Sending the response
	c.JSON(http.StatusOK, transactions)
}

//Function for fetching all transactions in a instance
func GetAllTransactionsHandler(c *gin.Context) {

	// Extract the authenticated user's ID and role from context.
	_, authUserRole, isAuthenticated := utils.GetAuthenticatedUser(c)
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or role not found"})
		return
	}

	// Ensure the authenticated user is an admin.
	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized. Admin access required"})
		return
	}

	// Database connection
	db := db.GetConnection(c)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection issue"})
		return
	}

	// Fetching all transactions directly within the handler.
	query := "SELECT * FROM transactions"
	var transactions []models.TransactionDisplay
	if err := db.Select(&transactions, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions", "details": err.Error()})
		return
	}

	// Handling empty transactions
	if len(transactions) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No transactions found"})
		return
	}

	// Sending the response
	c.JSON(http.StatusOK, transactions)
}
