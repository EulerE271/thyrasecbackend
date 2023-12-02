// holdingsRepository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	// Other imports, like uuid package
)

type AccountBalanceRepository struct {
	db *sql.DB
}

func NewAccountBalanceRepository(db *sql.DB) *AccountBalanceRepository {
	return &AccountBalanceRepository{db: db}
}

// GetTotalValue fetches the total value, cash, assets value, and available cash for a user.
func (r *AccountBalanceRepository) GetAggregatedValue(ctx context.Context, userId uuid.UUID) (TotalValue, error) {
	var totalValue TotalValue

	// SQL Query
	query := `
    SELECT 
    SUM(ac.account_balance) AS total_cash,
    COALESCE(SUM(h.quantity * a.current_price), 0) AS assets_value,
    SUM(ac.account_balance) + COALESCE(SUM(h.quantity * a.current_price), 0) AS total_value,
    SUM(ac.available_cash) AS available_cash
FROM 
    thyrasec.accounts ac
LEFT JOIN 
    thyrasec.holdings h ON ac.id = h.account_id
LEFT JOIN 
    thyrasec.assets a ON h.asset_id = a.id
WHERE 
    ac.account_holder_id = $1
GROUP BY 
    ac.account_holder_id;
    `
	fmt.Println("userId: %v", userId)
	// Execute the query
	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&totalValue.TotalCash,
		&totalValue.AssetValue,
		&totalValue.TotalValue,
		&totalValue.AvailableCash,
	)
	if err != nil {
		return TotalValue{}, err // Return an empty TotalValue struct and the error
	}

	return totalValue, nil
}

// TotalValue struct to hold the aggregated values
type TotalValue struct {
	TotalValue    float64
	TotalCash     float64
	AssetValue    float64
	AvailableCash float64
}

func (r *AccountBalanceRepository) GetSpecificAccountValue(ctx context.Context, accountId uuid.UUID) (AccountValue, error) {
	var account AccountValue

	// Adjusted SQL Query to fetch details of a specific account
	query := `
    SELECT 
        ac.id,
        ac.account_name,
        ac.account_balance,
        SUM(h.quantity * a.current_price) AS assets_value,
        ac.available_cash
    FROM 
        thyrasec.accounts ac
    LEFT JOIN 
        thyrasec.holdings h ON ac.id = h.account_id
    LEFT JOIN 
        thyrasec.assets a ON h.asset_id = a.id
    WHERE 
        ac.id = $1
    GROUP BY 
        ac.id;
    `

	err := r.db.QueryRowContext(ctx, query, accountId).Scan(
		&account.AccountID,
		&account.AccountName,
		&account.AccountBalance,
		&account.AssetValue,
		&account.AvailableCash,
	)
	if err != nil {
		return AccountValue{}, err
	}

	return account, nil
}

// AccountValue struct to hold individual account values
type AccountValue struct {
	AccountID      uuid.UUID
	AccountName    string
	AccountBalance float64
	AssetValue     float64
	AvailableCash  float64
}
