package utils

import (
	"context"
	"database/sql" // Import the standard sql package for sql.ErrNoRows
	"fmt"
	"log"
	"thyra/internal/common/db" // Import your custom db package

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func GetHouseAccount(c *gin.Context) (string, error) {
	// Use GetDB to get the initialized database connection
	database := db.GetDB()
	if database == nil {
		return "", fmt.Errorf("database connection is not initialized")
	}

	// Extract context from gin.Context, use context.Background() as fallback
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Struct to hold query result
	var result struct {
		HouseAccountTypeID string `db:"id"`
	}

	// Query to get the ID of the house account type
	accountTypeQuery := `SELECT id FROM thyrasec.account_types WHERE account_type_name = 'House'`
	err := database.GetContext(ctx, &result, accountTypeQuery) // Use 'database', not 'db'
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("House account type not found")
			return "", fmt.Errorf("house account type not found")
		}
		log.Printf("Error querying account type database: %v", err)
		return "", fmt.Errorf("error querying account type database: %v", err)
	}

	// Query to get the house account using the account type ID
	var houseAccountID string
	accountQuery := `SELECT id FROM thyrasec.accounts WHERE account_type = $1`
	err = database.GetContext(ctx, &houseAccountID, accountQuery, result.HouseAccountTypeID) // Use 'database', not 'db'
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No house account found")
			return "", fmt.Errorf("no house account found")
		}
		log.Printf("Error querying accounts database: %v", err)
		return "", fmt.Errorf("error querying accounts database: %v", err)
	}

	log.Printf("House account ID: %s", houseAccountID)

	return houseAccountID, nil
}
