package handlers

import (
	"net/http"
	"thyra/internal/orders/models"
	"thyra/internal/orders/services"

	"github.com/gin-gonic/gin"
)

type SettlementHandler struct {
	SetlementService *services.SettlementService
}

func NewSettlementHandler(settlementService *services.SettlementService) *SettlementHandler {
	return &SettlementHandler{
		SetlementService: settlementService,
	}
}

func (h *SettlementHandler) SettlementSellHandler(c *gin.Context) {
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

	orderID := c.Param("orderId")

	settlementRequest := models.SettlementRequest{}
	if err := c.BindJSON(&settlementRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	err := h.SetlementService.SellOrder(c, orderID, userIDStr, settlementRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to settle order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order settled successfully"})
}
