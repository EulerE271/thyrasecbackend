package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"thyra/internal/assets/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

/* Inserts a asset reservation when selling assets */
func InsertReservation(tx *sqlx.Tx, order models.Order, reservedUntil time.Time) error {
	reservation := struct {
		OrderID       uuid.UUID `db:"order_id"`
		AccountID     uuid.UUID `db:"account_id"`
		AssetID       uuid.UUID `db:"asset_id"`
		Quantity      int       `db:"quantity"`
		ReservedUntil time.Time `db:"reserved_until"`
		Status        string    `db:"status"`
	}{
		OrderID:       order.ID,
		AccountID:     order.AccountID,
		AssetID:       order.AssetID,
		Quantity:      int(order.Quantity),
		ReservedUntil: reservedUntil,
		Status:        "reserved",
	}

	query := `INSERT INTO thyrasec.reservations (order_id, account_id, asset_id, quantity, reserved_until, status, created_at, updated_at) 
              VALUES (:order_id, :account_id, :asset_id, :quantity, :reserved_until, :status, NOW(), NOW())`

	_, err := tx.NamedExec(query, reservation)
	return err
}

/* inserts cash reservation when pruchasing assets */
func InsertCashReservation(tx *sqlx.Tx, order models.Order, reservedUntil time.Time) error {
	cashReservation := struct {
		OrderID       uuid.UUID       `db:"order_id"`
		AccountID     uuid.UUID       `db:"account_id"`
		Amount        decimal.Decimal `db:"amount"`
		ReservedUntil time.Time       `db:"reserved_until"`
		Status        string          `db:"status"`
	}{
		OrderID:       order.ID,
		AccountID:     order.AccountID,
		Amount:        decimal.NewFromFloat(order.TotalAmount),
		ReservedUntil: reservedUntil,
		Status:        "reserved",
	}

	query := `INSERT INTO thyrasec.cash_reservations (order_id, account_id, amount, reserved_until, status, created_at, updated_at) 
              VALUES (:order_id, :account_id, :amount, :reserved_until, :status, NOW(), NOW())`

	_, err := tx.NamedExec(query, cashReservation)
	return err
}

/*Updates the available_cash and reserved_cash column in accounts table */
func ReserveCash(tx *sqlx.Tx, accountID uuid.UUID, amount decimal.Decimal) error {
	updateQuery := `UPDATE thyrasec.accounts SET 
                    available_cash = available_cash - $1, 
                    reserved_cash = reserved_cash + $1 
                    WHERE id = $2 AND available_cash >= $1`

	_, err := tx.Exec(updateQuery, amount, accountID)
	return err
}

/* Deducts the holding when selling an asset */
func DeductHoldings(tx *sqlx.Tx, accountID, assetID uuid.UUID, quantity int) error {
	updateQuery := `UPDATE thyrasec.holdings 
                    SET quantity = quantity - :quantity 
                    WHERE account_id = :account_id AND asset_id = :asset_id AND quantity >= :quantity`

	params := map[string]interface{}{
		"quantity":   quantity,
		"account_id": accountID,
		"asset_id":   assetID,
	}

	result, err := tx.NamedExec(updateQuery, params)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no holdings updated, possibly due to insufficient quantity")
	}

	return nil
}

func GetOrder(db *sqlx.DB, orderID string) (*models.Order, error) {
	var order models.Order
	query := "SELECT * FROM thyrasec.orders WHERE id = $1"
	if err := db.Get(&order, query, orderID); err != nil {
		return nil, err
	}
	return &order, nil
}

func GetOrderType(db *sqlx.DB, trtShortName string) (uuid.UUID, error) {

	var trtUUID uuid.UUID
	query := "SELECT type_id FROM transactions_types WHERE trt_short_name = $1"
	if err := db.Get(&trtUUID, query, trtShortName); err != nil {
		return uuid.Nil, err
	}

	return trtUUID, nil
}

func UpdateOrderStatus(tx *sqlx.Tx, orderID string, status models.OrderStatusType) error {
	query := "UPDATE thyrasec.orders SET status = $1 WHERE id = $2"
	_, err := tx.Exec(query, status, orderID)
	return err
}

/*Checks whether there is enough holdings to sell an asset */
func CheckHoldings(tx *sqlx.Tx, accountID, assetID uuid.UUID) (int, error) {
	var currentQuantity int
	query := `SELECT quantity FROM thyrasec.holdings WHERE account_id = $1 AND asset_id = $2`
	err := tx.Get(&currentQuantity, query, accountID, assetID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No holdings found for the account and asset
			return 0, errors.New("no holdings found for the specified account and asset")
		}
		return 0, err
	}
	return currentQuantity, nil
}

/* Checks whether there is enoguht cash when buying an asset */
func CheckAvailableCash(tx *sqlx.Tx, accountID uuid.UUID, amount decimal.Decimal) (bool, error) {
	var availableCash decimal.Decimal
	query := `SELECT available_cash FROM thyrasec.accounts WHERE id = $1`
	err := tx.Get(&availableCash, query, accountID)
	if err != nil {
		return false, err
	}
	return availableCash.GreaterThanOrEqual(amount), nil
}

func ReleaseReservation(db *sqlx.DB, orderID string) error {
	var accountID uuid.UUID
	var amount decimal.Decimal

	// Prepare the query
	quantityQuery := `SELECT account_id, amount FROM thyrasec.cash_reservations WHERE order_id = $1`

	// Execute the query and scan the results into the variables
	err := db.QueryRow(quantityQuery, orderID).Scan(&accountID, &amount)
	if err != nil {
		return fmt.Errorf("failed to get reservation details: %w", err)
	}

	// Update account reserved cash
	updateAccountQuery := `UPDATE thyrasec.accounts SET reserved_cash = reserved_cash - $1 WHERE id = $2`
	_, err = db.Exec(updateAccountQuery, amount, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account reserved cash: %w", err)
	}

	// Update reservation status
	updateReservationQuery := `UPDATE thyrasec.cash_reservations SET status = 'inactive' WHERE order_id = $1`
	_, err = db.Exec(updateReservationQuery, orderID)
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	return nil
}

func UpdateAccountBalance(db *sqlx.DB, accountID uuid.UUID, balanceChange float64) error {
	// Start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	fmt.Sprintf("Updating inside function: %v, %w", accountID, balanceChange)
	// Define the query to update the account
	query := `UPDATE thyrasec.accounts
	SET account_balance = account_balance - $1,
		reserved_cash = reserved_cash - $1,
		available_cash = available_cash - $1
	WHERE id = $2`

	// Execute the query
	_, err = tx.Exec(query, balanceChange, accountID)
	if err != nil {
		tx.Rollback() // Roll back the transaction on error
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertHolding(db *sqlx.DB, holding models.Holding) error {
	query := `INSERT INTO thyrasec.holdings (id, account_id, asset_id, quantity)
			  VALUES (:id, :account_id, :asset_id, :quantity)`

	_, err := db.NamedExec(query, holding)
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrder(db *sqlx.DB, orderID string, settledQuantity float64, settledAmount float64, status string, tradeDate *time.Time, settlementDate *time.Time, comment string) error {
	query := `
        UPDATE thyrasec.orders
        SET settledQuantity = $1, settledAmount = $2, status = $3, trade_date = $4, settlement_date = $5, comment = $6, updated_at = NOW()
        WHERE id = $7	
        `

	_, err := db.Exec(query, settledQuantity, settledAmount, status, tradeDate, settlementDate, comment, orderID)
	return err
}
