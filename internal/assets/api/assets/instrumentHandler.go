package handlers

import (
	"net/http"
	"thyra/internal/assets/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CreateInstrument creates a new instrument.
func CreateInstrument(c *gin.Context) {
	// Extract the authenticated user's ID and role from context.
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to create an instrument"})
		return
	}

	var instrument models.Instrument
	if err := c.ShouldBindJSON(&instrument); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	// Check if the provided asset_type_id exists in the asset_types table
	var count int
	err := sqlxDB.Get(&count, "SELECT COUNT(*) FROM thyrasec.asset_types WHERE id = $1", instrument.AssetTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify asset_type_id", "details": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset_type_id"})
		return
	}

	query := `
    INSERT INTO assets (instrument_name, isin, ticker, exchange, currency, instrument_type, current_price, volume, country, sector, asset_type_id) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING id, instrument_name, isin, ticker, exchange, currency, instrument_type, current_price, volume, country, sector, asset_type_id
`
	row := sqlxDB.QueryRowx(query, instrument.InstrumentName, instrument.ISIN, instrument.Ticker, instrument.Exchange, instrument.Currency, instrument.InstrumentType, instrument.CurrentPrice, instrument.Volume, instrument.Country, instrument.Sector, instrument.AssetTypeID)
	if err = row.StructScan(&instrument); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create instrument", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instrument created successfully", "instrument": instrument})

}

// GetAllInstruments fetches all instruments.
func GetAllInstruments(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to fetch all instruments"})
		return
	}

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	query := "SELECT * FROM assets"
	var instruments []models.Instrument
	if err := sqlxDB.Select(&instruments, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch instruments", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instruments)
}

func GetAllAssetTypes(c *gin.Context) {

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to fetch all instruments"})
		return
	}

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	query := "SELECT * FROM asset_types"
	var assets []models.Asset
	if err := sqlxDB.Select(&assets, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch asset types", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assets)

}
