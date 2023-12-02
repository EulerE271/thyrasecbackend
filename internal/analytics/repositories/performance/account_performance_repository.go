package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type AccountPerformanceRepository struct {
	db *sql.DB
}

func NewAccountPerformanceRepository(db *sql.DB) *AccountPerformanceRepository {
	return &AccountPerformanceRepository{
		db: db,
	}
}

type ValueChange struct {
	StartValue       float64
	EndValue         float64
	Change           float64
	PercentualChange float64
	StartDate        time.Time
	EndDate          time.Time
	Snapshots        []SnapshotValue // Add this field to your struct
}

type SnapshotValue struct {
	Date  time.Time
	Value float64
}

func (r *AccountPerformanceRepository) GetAccountPerformanceChange(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) (ValueChange, error) {
	snapshots, err := r.fetchAccountSnapshots(ctx, accountID, startDate, endDate)
	if err != nil {
		return ValueChange{}, err
	}
	transactions, err := r.fetchTransactions(ctx, accountID, startDate, endDate)
	if err != nil {
		return ValueChange{}, err
	}

	totalCashFlow := calculateCashFlows(transactions)
	// Use the first and last snapshots for start and end values
	startValue := snapshots[0].Value
	endValue := snapshots[len(snapshots)-1].Value

	// Adjust end value for cash flows
	adjustedEndValue := endValue - totalCashFlow

	// Calculate value change using adjusted end value
	valueChange := adjustedEndValue - startValue
	percentualChange := 0.0
	if startValue != 0 {
		percentualChange = (valueChange / startValue) * 100
	}

	return ValueChange{
		StartValue:       startValue,
		EndValue:         adjustedEndValue, // Use adjusted end value here
		Change:           valueChange,
		PercentualChange: percentualChange,
		StartDate:        snapshots[0].Date,
		EndDate:          snapshots[len(snapshots)-1].Date,
		Snapshots:        snapshots,
	}, nil
}

func (r *AccountPerformanceRepository) calculateTotalMarketValue(ctx context.Context, accountID uuid.UUID, date time.Time) (float64, time.Time, error) {
	var totalMarketValue float64
	var actualDate time.Time

	// Adjusted query to handle dates before the first recorded holding
	query := `
	WITH ranked_snapshots AS (
		SELECT hs.quantity, ap.price, hs.snapshot_date,
		ROW_NUMBER() OVER (ORDER BY ABS(hs.snapshot_date - $2)) as rn
		FROM thyrasec.holdings_snapshots hs
		JOIN thyrasec.asset_prices ap ON hs.asset_id = ap.asset_id
		WHERE hs.account_id = $1
	)
	SELECT quantity, price, snapshot_date FROM ranked_snapshots WHERE rn = 1;
	`

	row := r.db.QueryRowContext(ctx, query, accountID, date)
	var quantity, price float64
	if err := row.Scan(&quantity, &price, &actualDate); err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where there is no data
			return 0, time.Time{}, fmt.Errorf("no data available for account %s", accountID)
		}
		return 0, time.Time{}, err
	}

	totalMarketValue = quantity * price
	return totalMarketValue, actualDate, nil
}

//Returns the average performance of all of a users account
//Returns the average performance of all of a users account
func (r *AccountPerformanceRepository) GetUserPerformanceChange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (ValueChange, error) {
	accountIDs, err := r.fetchUserAccountIDs(ctx, userID)
	if err != nil {
		return ValueChange{}, err
	}

	snapshotMap := make(map[time.Time]float64)
	snapshotCount := make(map[time.Time]int)

	for _, accountID := range accountIDs {
		accountSnapshots, err := r.fetchAggregatedAccountSnapshots(ctx, accountID, startDate, endDate)
		if err != nil {
			return ValueChange{}, err
		}

		// Aggregate snapshot values by date
		for _, snapshot := range accountSnapshots {
			snapshotMap[snapshot.Date] += snapshot.Value
			snapshotCount[snapshot.Date]++
		}
	}

	// Check if there are any snapshots
	if len(snapshotMap) == 0 {
		return ValueChange{}, fmt.Errorf("no data available for user %s", userID)
	}

	// Calculate average snapshot values and sort them by date
	var snapshots []SnapshotValue
	for date, totalValue := range snapshotMap {
		avgValue := totalValue / float64(snapshotCount[date])
		snapshots = append(snapshots, SnapshotValue{Date: date, Value: avgValue})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Date.Before(snapshots[j].Date)
	})

	startValue := snapshots[0].Value
	endValue := snapshots[len(snapshots)-1].Value

	valueChange := endValue - startValue
	percentualChange := 0.0
	if startValue != 0 {
		percentualChange = (valueChange / startValue) * 100
	}

	return ValueChange{
		StartValue:       startValue,
		EndValue:         endValue,
		Change:           valueChange,
		PercentualChange: percentualChange,
		StartDate:        startDate,
		EndDate:          endDate,
		Snapshots:        snapshots,
	}, nil
}

func (r *AccountPerformanceRepository) fetchAggregatedAccountSnapshots(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) ([]SnapshotValue, error) {
	// Query to fetch snapshot data between startDate and endDate
	var snapshots []SnapshotValue
	query := `SELECT hs.snapshot_date, hs.quantity * ap.price as value
	FROM thyrasec.holdings_snapshots hs
	JOIN thyrasec.asset_prices ap ON hs.asset_id = ap.asset_id
	AND ap.price_date = (
		SELECT MAX(price_date)
		FROM thyrasec.asset_prices
		WHERE asset_id = hs.asset_id AND price_date <= hs.snapshot_date
	)
	WHERE hs.account_id = $1 AND hs.snapshot_date BETWEEN $2 AND $3
	`
	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snapshot SnapshotValue
		if err := rows.Scan(&snapshot.Date, &snapshot.Value); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}

//Fetches holding snapshot
func (r *AccountPerformanceRepository) fetchAccountSnapshots(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) ([]SnapshotValue, error) {
	var snapshots []SnapshotValue

	// Updated query to fetch the most recent price for each snapshot date
	query := `
	SELECT hs.snapshot_date, hs.quantity * ap.price as value
	FROM thyrasec.holdings_snapshots hs
	JOIN thyrasec.asset_prices ap ON hs.asset_id = ap.asset_id
	AND ap.price_date = (
	    SELECT MAX(ap_inner.price_date)
	    FROM thyrasec.asset_prices ap_inner
	    WHERE ap_inner.asset_id = hs.asset_id AND ap_inner.price_date <= hs.snapshot_date
	)
	WHERE hs.account_id = $1 AND hs.snapshot_date BETWEEN $2 AND $3
	ORDER BY hs.snapshot_date
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snapshot SnapshotValue
		if err := rows.Scan(&snapshot.Date, &snapshot.Value); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}

//Returns the ID of all accounts
func (r *AccountPerformanceRepository) fetchUserAccountIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var accountIDs []uuid.UUID

	query := "SELECT id FROM thyrasec.accounts WHERE account_holder_id = $1"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var accountID uuid.UUID
		if err := rows.Scan(&accountID); err != nil {
			return nil, err
		}
		accountIDs = append(accountIDs, accountID)
	}

	return accountIDs, nil
}

type Transaction struct {
	Id                 uuid.UUID `json:"id" db:"id"`
	TransactionOwnerId uuid.UUID `json:"transaction_owner_id" db:"transaction_owner_id"`
	AccountOwnerId     uuid.UUID `json:"account_owner_id" db:"account_owner_id"`
	Type               uuid.UUID `json:"type" db:"type"`
	Asset1Id           uuid.UUID `json:"asset1_id" db:"asset1_id"`                 // Nullable field
	Asset2Id           uuid.UUID `json:"asset2_id" db:"asset2_id"`                 // Nullable field
	AccountAsset1Id    uuid.UUID `json:"account_asset1_id" db:"account_asset1_id"` // Nullable field
	AccountAsset2Id    uuid.UUID `json:"account_asset2_id" db:"account_asset2_id"` // Nullable field
	AmountAsset1       *float64  `json:"amount_asset1" db:"amount_asset1"`         // Nullable field
	AmountAsset2       *float64  `json:"amount_asset2" db:"amount_asset2"`         // Nullable field
	CreatedById        uuid.UUID `json:"created_by_id" db:"created_by_id"`
	UpdatedById        uuid.UUID `json:"updated_by_id" db:"updated_by_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	Corrected          bool      `json:"corrected" db:"corrected"`
	Canceled           bool      `json:"canceled" db:"canceled"`
	Comment            *string   `json:"comment" db:"comment"` // Nullable field
	Asset2_currency    uuid.UUID `json:"asset2_currency" db:"asset2_currency"`
	Asset2_price       uuid.UUID `json:"asset2_price" db:"asset2_price"`
	OrderNumber        string    `json:"order_no" db:"order_no"`
	BusinessEvent      uuid.UUID `json:"business_event" db:"business_event"`
	Trade_date         time.Time `json:"trade_date" db:"trade_date"`
	Settlement_date    time.Time `json:"settlement_date" db:"settlement_date"`
}

func (r *AccountPerformanceRepository) fetchTransactions(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) ([]Transaction, error) {
	var transactions []Transaction

	// Query to fetch transactions
	query := `
    SELECT id, type, amount_asset1, amount_asset2, trade_date
    FROM thyrasec.transactions
    WHERE account_owner_id = $1 AND trade_date BETWEEN $2 AND $3
    `

	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.Id, &transaction.Type, &transaction.AmountAsset1, &transaction.AmountAsset2, &transaction.Trade_date); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func calculateCashFlows(transactions []Transaction) float64 {
	var totalCashFlow float64
	for _, transaction := range transactions {
		// Check for deposit
		if transaction.AmountAsset1 != nil && *transaction.AmountAsset1 > 0 {
			totalCashFlow += *transaction.AmountAsset1 // Adding deposit to cash flow
		}

		// Check for withdrawal
		if transaction.AmountAsset2 != nil && *transaction.AmountAsset2 > 0 {
			totalCashFlow -= *transaction.AmountAsset2 // Subtracting withdrawal from cash flow
		}
	}
	return totalCashFlow
}
