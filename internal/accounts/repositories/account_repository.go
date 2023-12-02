package repository

import (
	"context"
	"database/sql"
	"log" // Added import for logging

	"github.com/google/uuid"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetHouseAccount(ctx context.Context) (uuid.UUID, error) {
	query := `
        SELECT a.id
        FROM accounts a
        INNER JOIN account_types at ON a.account_type = at.id
        WHERE at.account_type_name = $1;
    `

	var accountID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, "House").Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no house account is found
			log.Printf("No house account found: %v", err)
			return uuid.Nil, err // Return an empty UUID and the error
		}
		// Handle other errors
		log.Printf("Error querying house account: %v", err)
		return uuid.Nil, err // Return an empty UUID and the error
	}

	return accountID, nil // Return the result and nil error
}
