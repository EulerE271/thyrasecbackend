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

//Calculates the performance of a single account
func (r *AccountPerformanceRepository) GetAccountPerformanceChange(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) (ValueChange, error) {
	snapshots, err := r.fetchAccountSnapshots(ctx, accountID, startDate, endDate)
	if err != nil {
		return ValueChange{}, err
	}

	// Check if there are any snapshots
	if len(snapshots) == 0 {
		return ValueChange{}, fmt.Errorf("no data available for account %s", accountID)
	}

	// Use the first and last snapshots for start and end values
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
