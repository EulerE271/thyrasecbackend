package handlers

import (
	"net/http"
	"thyra/internal/common/db"
	"thyra/internal/transactions/repositories"
	"thyra/internal/transactions/services"
	"thyra/internal/transactions/utils"

	"github.com/gin-gonic/gin"
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
	repo := repositories.NewTransactionRepository(database)
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
	database := db.GetConnection(c)
	repo := repositories.NewTransactionRepository(database)
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
