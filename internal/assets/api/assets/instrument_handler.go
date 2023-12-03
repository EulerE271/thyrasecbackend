package handlers

import (
	"net/http"
	"thyra/internal/assets/models"
	"thyra/internal/assets/services"

	"github.com/gin-gonic/gin"
)

type AssetsHandler struct {
	service *services.AssetsService
}

func NewAssetsHandler(service *services.AssetsService) *AssetsHandler {
	return &AssetsHandler{service: service}
}

func (h *AssetsHandler) CreateInstrument(c *gin.Context) {
	var instrument models.Instrument
	if err := c.ShouldBindJSON(&instrument); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole, _ := c.Get("userType")
	createdInstrument, err := h.service.CreateInstrument(c.Request.Context(), instrument, userRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instrument created successfully", "instrument": createdInstrument})
}

func (h *AssetsHandler) GetAllInstruments(c *gin.Context) {
	userRole, _ := c.Get("userType")

	instruments, err := h.service.GetAllInstruments(c.Request.Context(), userRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instruments)
}

func (h *AssetsHandler) GetAllAssetTypes(c *gin.Context) {
	userRole, _ := c.Get("userType")

	assets, err := h.service.GetAllAssetTypes(c.Request.Context(), userRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assets)
}
