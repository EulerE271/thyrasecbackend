package utils

import (
	"log"
	"thyra/internal/transactions/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func BindTransaction(c *gin.Context) (*models.Transaction, error) {
	var newTransaction models.Transaction
	err := c.ShouldBindWith(&newTransaction, binding.JSON)
	if err != nil {
		// This will help you see the exact reason for the binding failure
		log.Printf("BindTradeTransaction error: %v", err)
		return nil, err
	}
	return &newTransaction, err
}
