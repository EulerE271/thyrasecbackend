package handlers

import (
	"net/http"
	account "thyra/internal/accounts/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountValueHandler struct {
	service *account.AccountValueService
}

func NewAccountValueHandler(service *account.AccountValueService) *AccountValueHandler {
	return &AccountValueHandler{service: service}
}

// GetAggregatedValues handles requests for fetching aggregated values for a user.
func (h *AccountValueHandler) GetAggregatedValues(c *gin.Context) {
	// Extract userID from URL parameter instead of context
	userID := c.Param("userId")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in URL"})
		return
	}

	// Fetch aggregated values using the service
	totalValue, err := h.service.GetAggregatedAccountValue(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregated values", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, totalValue)
}

func (h *AccountValueHandler) GetSpecificAccountValue(c *gin.Context) {
	// Extract accountId from the request
	accountId, err := uuid.Parse(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Fetch specific account values using the service
	account, err := h.service.GetSpecificAccountValue(c.Request.Context(), accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account values", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}
