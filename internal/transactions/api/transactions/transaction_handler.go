package handlers

import (
	"database/sql"
	"net/http"
	"thyra/internal/common/db"
	"thyra/internal/transactions/models"
	"thyra/internal/transactions/repositories"
	"thyra/internal/transactions/services"
	"thyra/internal/transactions/utils"
	userutils "thyra/internal/users/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateDeposit(c *gin.Context) {
	database := db.GetConnection(c)
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error database": "Internal Server Error"})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error userIdInterface": "UserID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error userIDStr": "UserID is not a string"})
		return
	}

	// Parse the userID string as UUID (only if needed later in the code)
	/* userUUID, err := uuid.Parse(userIDStr)
	   if err != nil {
	       c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is not a valid UUID", "details": err.Error()})
	       return
	   } */

	newTransaction, err := utils.BindTransaction(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input in main transaction data", "error": err.Error()})
		return
	}

	// Create a new repository and service instance
	repo := repositories.NewTransactionRepository(&sqlx.Tx{})
	service := services.NewTransactionService(repo)

	// Use the service to create the deposit
	debitTransactionID, creditTransactionID, err := service.CreateDeposit(c, userIDStr, newTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error createdeposit": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "Parent debitTransactionID": debitTransactionID, "Parent CreditTransactionID": creditTransactionID})
}

func CreateWithdrawal(c *gin.Context) {

	// Extract user ID from the context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a string"})
		return
	}

	// Bind transaction data from the request body
	newTransaction, err := utils.BindTransaction(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input in main transaction data", "error": err.Error()})
		return
	}

	// Create a new service instance
	//tx := db.GetConnection(c)
	repo := repositories.NewTransactionRepository(&sqlx.Tx{})
	service := services.NewTransactionService(repo)

	// Use the service to create the withdrawal
	debitTransactionID, creditTransactionID, err := service.CreateWithdrawal(c, userIDStr, newTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":                    "Transaction created successfully",
		"Parent debitTransactionID":  debitTransactionID,
		"Parent CreditTransactionID": creditTransactionID,
	})
}

//Function for fetching all transactions for a specific user
func GetTransactionByUserHandler(c *gin.Context) {
	// Extracting user parameters
	targetUserID := c.Param("userId")
	authUserID, authUserRole, isAuthenticated := userutils.GetAuthenticatedUser(c)

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
        a.account_number
    FROM 
        thyrasec.transactions t
    LEFT JOIN 
        thyrasec.transactions_types tt ON t.type = tt.type_id
    LEFT JOIN
        thyrasec.accounts a ON t.transaction_owner_account_id = a.id
    WHERE 
        t.transaction_owner_account_id = $1;
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
	_, authUserRole, isAuthenticated := userutils.GetAuthenticatedUser(c)
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

func GetAssetID(c *gin.Context) {
	db := db.GetConnection(c)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection issue"})
		return
	}

	identifier := c.Query("identifier")

	query := `
        SELECT c.id, c.unified_asset_id, ua.asset_type_id
        FROM thyrasec.currencies c
        JOIN thyrasec.unified_assets ua ON c.unified_asset_id = ua.id
        WHERE ua.identifier = $1;
    `

	var assetID, unifiedAssetId, assetTypeId string
	err := db.QueryRow(query, identifier).Scan(&assetID, &unifiedAssetId, &assetTypeId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"asset_id": assetID, "unified_asset_id": unifiedAssetId, "asset_type_id": assetTypeId})
}

func GetTransactionTypesHandler(c *gin.Context) {
	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	var transactionTypes []models.TransactionType
	if err := sqlxDB.Select(&transactionTypes, "SELECT * FROM transactions_types"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction types", "details": err.Error()})
		return
	}

	if len(transactionTypes) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No transaction types found"})
		return
	}

	c.JSON(http.StatusOK, transactionTypes)
}
