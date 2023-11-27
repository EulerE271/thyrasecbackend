package services

import (
	"database/sql"
	"log"
	"thyra/internal/common/db"

	"github.com/google/uuid"
)

// GetHouseAccount returns the ID of the house account and an error, if any
func GetHouseAccount() uuid.UUID {
	// Get the database connection
	sqlxDB := db.GetDB()

	// Define the query to find the house account
	query := `
        SELECT a.id
        FROM accounts a
        INNER JOIN account_types at ON a.account_type = at.id
        WHERE at.account_type_name = $1;
    `

	var accountID uuid.UUID
	// Execute the query with "House" as the account type name
	err := sqlxDB.Get(&accountID, query, "House")
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no house account is found
			log.Printf("No house account found: %v", err)
		}
		// Handle other errors
		log.Printf("Error querying house account: %v", err)
	}

	return accountID
}
