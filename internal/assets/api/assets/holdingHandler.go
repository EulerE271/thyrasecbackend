package handlers

import (
	"log"
	"net/http"
	"thyra/internal/assets/services" // Replace with the actual path to your services

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Assuming you have a way to initialize this with a HoldingsService instance
type HoldingsHandler struct {
	service services.HoldingsService
}

// NewHoldingsHandler creates a new instance of HoldingsHandler
func NewHoldingsHandler(service services.HoldingsService) *HoldingsHandler {
	return &HoldingsHandler{service: service}
}

// GetAccountHoldings handles the request to get account holdings
func (h *HoldingsHandler) GetAccountHoldings(c *gin.Context) {
	accountIdStr := c.Param("accountId")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	holdings, err := h.service.GetAccountHoldings(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve holdings"})
		return
	}

	c.JSON(http.StatusOK, holdings)
}

// GetAssetDetails handles the request to get asset details
func (h *HoldingsHandler) GetAssetDetails(c *gin.Context) {
	assetIdStr := c.Param("assetId")
	assetId, err := uuid.Parse(assetIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}

	asset, err := h.service.GetAssetDetails(assetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset details"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *HoldingsHandler) GetAccountHoldingsWithDetails(c *gin.Context) {
	accountIdStr := c.Param("accountId")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	holdingsWithDetails, err := h.service.GetHoldingsWithAssetDetails(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messag": "Failed to retrieve holdings with details", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, holdingsWithDetails)
}

func (h *HoldingsHandler) GetCurrencyID(c *gin.Context) {
	currencyName := c.Query("currency")

	currencyID, err := h.service.GetCurrencyID(currencyName)
	if err != nil {
		log.Printf("Error fetching currency ID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch currency ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currencyID": currencyID,
	})
}
