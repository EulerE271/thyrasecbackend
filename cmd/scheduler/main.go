package main

import (
	"fmt"
	"log"
	"thyra/internal/common/db"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func main() {
	const maxRetries = 50
	const retryInterval = 10 * time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
		err = db.Initialize()
		if err == nil {
			break
		}

		log.Printf("Failed to initialize the database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to initialize the database after %d attempts: %v", maxRetries, err)
	}

	// Test database connection
	testDB := db.GetDB()
	if testDB == nil {
		log.Fatal("Database connection is nil after initialization")
	}

	var testVar int
	err = testDB.Get(&testVar, "SELECT 1")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}

	log.Printf("Database is active")
	performScheduledTasks(testDB)

	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			performScheduledTasks(testDB)
		}
	}
}

func performScheduledTasks(db *sqlx.DB) {
	accountIDs, err := fetchAllAccounts(db)
	if err != nil {
		log.Printf("Error fetching accounts: %v", err)
		return
	}

	snapshotDate := time.Now().AddDate(0, 0, -1)
	for _, accountID := range accountIDs {
		tx, err := db.Beginx()
		if err != nil {
			log.Printf("Failed to start a transaction: %v", err)
			continue
		}

		err = CalculateAndStoreHoldings(tx, accountID, snapshotDate)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to calculate and store performance for account %s: %v", accountID, err)
		} else {
			tx.Commit()
		}
	}
}

type HoldingSnapshot struct {
	SnapshotID   uuid.UUID `db:"snapshotid"` // Assuming you want to handle this in Go
	AccountID    uuid.UUID `db:"account_id"`
	AssetID      uuid.UUID `db:"asset_id"`
	Quantity     float64   `db:"quantity"`
	SnapshotDate time.Time `db:"snapshot_date"`
}
type Holding struct {
	AssetID  uuid.UUID `db:"asset_id"`
	Quantity float64   `db:"quantity"`
}

func CalculateAndStoreHoldings(tx *sqlx.Tx, accountID uuid.UUID, snapshotDate time.Time) error {
	holdings, err := fetchHoldings(tx, accountID)
	if err != nil {
		return fmt.Errorf("error fetching holdings for account %s: %v", accountID, err)
	}

	for _, holding := range holdings {
		snapshot := HoldingSnapshot{
			AccountID:    accountID,
			AssetID:      holding.AssetID,
			Quantity:     holding.Quantity,
			SnapshotDate: snapshotDate,
		}

		query := "INSERT INTO thyrasec.holdings_snapshots (account_id, asset_id, quantity, snapshot_date) VALUES ($1, $2, $3, $4)"
		_, err = tx.Exec(query, snapshot.AccountID, snapshot.AssetID, snapshot.Quantity, snapshot.SnapshotDate)
		if err != nil {
			// Handle the error gracefully, log it, and continue to the next iteration
			log.Printf("Error inserting snapshot for asset %s in account %s: %v", holding.AssetID, accountID, err)
		}
	}

	return nil
}

func fetchHoldings(tx *sqlx.Tx, accountID uuid.UUID) ([]Holding, error) {
	var holdings []Holding

	query := "SELECT asset_id, quantity FROM thyrasec.holdings WHERE account_id = $1"
	err := tx.Select(&holdings, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("error fetching holdings for account %s: %v", accountID, err)
	}

	return holdings, nil
}

func fetchAllAccounts(db *sqlx.DB) ([]uuid.UUID, error) {
	var accountIDs []uuid.UUID

	query := "SELECT id FROM thyrasec.accounts"
	err := db.Select(&accountIDs, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching account IDs: %v", err)
	}

	return accountIDs, nil
}

/*func CalculateAndStorePerformance(tx *sqlx.Tx, accountID uuid.UUID, snapshotDate time.Time) error {
	totalValue, err := calculateTotalValue(tx, accountID, snapshotDate) // pass 'tx' as the first argument
	if err != nil {
		return fmt.Errorf("error calculating total value: %v", err)
	}

	Snapshot := AccountSnapshot{
		AccountID:    accountID,
		SnapshotDate: snapshotDate,
		TotalValue:   totalValue,
	}

	query := `INSERT INTO thyrasec.account_snapshots (account_id, snapshot_date, total_value) VALUES (:account_id, :snapshot_date, :total_value)`
	_, err = tx.NamedExec(query, Snapshot)
	if err != nil {
		return fmt.Errorf("error storing account snapshot: %v", err)
	}

	return nil
}
func calculateTotalValue(tx *sqlx.Tx, accountID uuid.UUID, snapshotDate time.Time) (float64, error) {
	var totalValue float64

	// Query to calculate the total value of holdings based on the snapshot date
	query := `
    SELECT COALESCE(SUM(h.quantity * p.current_price), 0)
    FROM thyrasec.holdings AS h
    INNER JOIN thyrasec.assets AS p ON h.asset_id = p.id
    WHERE h.account_id = $1`

	err := tx.Get(&totalValue, query, accountID)
	if err != nil {
		return 0, fmt.Errorf("error calculating total value: %v", err)
	}

	return totalValue, nil
}*/
