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
	startValue, actualStartDate, err := r.calculateTotalMarketValue(ctx, accountID, startDate)
	if err != nil {
		return ValueChange{}, err
	}

	endValue, actualEndDate, err := r.calculateTotalMarketValue(ctx, accountID, endDate)
	if err != nil {
		return ValueChange{}, err
	}

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
		StartDate:        actualStartDate,
		EndDate:          actualEndDate,
	}, nil
}

func (r *AccountPerformanceRepository) calculateTotalMarketValue(ctx context.Context, accountID uuid.UUID, date time.Time) (float64, time.Time, error) {
	var totalMarketValue float64
	var actualDate time.Time

	query := `
    SELECT hs.quantity, ap.price, hs.snapshot_date
    FROM thyrasec.holdings_snapshots hs
    JOIN thyrasec.asset_prices ap ON hs.asset_id = ap.asset_id
    WHERE hs.account_id = $1 AND hs.snapshot_date <= $2 AND ap.price_date <= $2
    ORDER BY hs.snapshot_date DESC, ap.price_date DESC
    LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, accountID, date)
	var quantity, price float64
	if err := row.Scan(&quantity, &price, &actualDate); err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where there is no data even for previous dates
			return 0, time.Time{}, fmt.Errorf("no data available for account %s up to date %s", accountID, date)
		}
		return 0, time.Time{}, err
	}

	totalMarketValue = quantity * price
	return totalMarketValue, actualDate, nil
}

func (r *AccountPerformanceRepository) GetUserPerformanceChange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (ValueChange, error) {
	accountIDs, err := r.fetchUserAccountIDs(ctx, userID)
	if err != nil {
		return ValueChange{}, err
	}

	var totalStartValue, totalEndValue float64
	snapshotMap := make(map[time.Time]float64)
	snapshotCount := make(map[time.Time]int)

	for _, accountID := range accountIDs {
		// ... existing code to calculate start/end values ...

		accountSnapshots, err := r.fetchAccountSnapshots(ctx, accountID, startDate, endDate)
		if err != nil {
			return ValueChange{}, err
		}

		// Aggregate snapshot values by date
		for _, snapshot := range accountSnapshots {
			snapshotMap[snapshot.Date] += snapshot.Value
			snapshotCount[snapshot.Date]++
		}
	}

	// Calculate average snapshot values
	var snapshots []SnapshotValue
	for date, totalValue := range snapshotMap {
		count := snapshotCount[date] // Get the count for this date
		if count == 0 {
			continue // Skip dates with zero count to avoid division by zero
		}
		avgValue := totalValue / float64(count)
		snapshots = append(snapshots, SnapshotValue{Date: date, Value: avgValue})
	}

	// Sort the combined snapshots by date
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Date.Before(snapshots[j].Date)
	})

	valueChange := totalEndValue - totalStartValue
	percentualChange := 0.0
	if totalStartValue != 0 {
		percentualChange = (valueChange / totalStartValue) * 100
	}

	return ValueChange{
		StartValue:       totalStartValue,
		EndValue:         totalEndValue,
		Change:           valueChange,
		PercentualChange: percentualChange,
		StartDate:        startDate,
		EndDate:          endDate,
		Snapshots:        snapshots,
	}, nil
}

func (r *AccountPerformanceRepository) fetchAccountSnapshots(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) ([]SnapshotValue, error) {
	// Query to fetch snapshot data between startDate and endDate
	var snapshots []SnapshotValue
	query := `SELECT snapshot_date, quantity * price as value
	FROM thyrasec.holdings_snapshots hs
	JOIN thyrasec.asset_prices ap ON hs.asset_id = ap.asset_id
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
