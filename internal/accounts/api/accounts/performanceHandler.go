package handlers

import (
	"net/http"
	"thyra/internal/accounts/services" // Replace with your actual package path
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountPerformanceHandler struct {
	accountService *services.AccountPerformanceService
}

func NewAccountPerformanceHandler(accountService *services.AccountPerformanceService) *AccountPerformanceHandler {
	return &AccountPerformanceHandler{
		accountService: accountService,
	}
}

func (h *AccountPerformanceHandler) GetAccountPerformanceChange(c *gin.Context) {
	// Extract account ID from request
	accountIDStr := c.Param("accountId")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	// Parse startDate and endDate from string to time.Time
	layout := "2006-01-02" // This is the Go layout for date format, adjust if your format is different
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// Now call the service with parsed time.Time values
	valueChange, err := h.accountService.GetAccountValueChange(c, accountID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, valueChange)
}

func (h *AccountPerformanceHandler) GetUserPerformanceChange(c *gin.Context) {
	userIdStr := c.Param("userId")
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse startDate and endDate from string to time.Time
	layout := "2006-01-02" // This is the Go layout for date format, adjust if your format is different
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}
	valueChange, err := h.accountService.GetUserValueChange(c, userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, valueChange)
}
