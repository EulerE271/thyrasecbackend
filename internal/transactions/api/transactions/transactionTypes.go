package handlers

import (
	"net/http"
	"thyra/internal/transactions/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx" // Import the sqlx package
)

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
