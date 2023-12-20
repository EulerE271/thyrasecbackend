package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"thyra/internal/orders/models"
	positionmodels "thyra/internal/positions/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type OrdersRepository struct {
	db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) *OrdersRepository {
	return &OrdersRepository{db: db}
}

/* Inserts a asset reservation when selling assets */
func (r *OrdersRepository) InsertReservation(tx *sqlx.Tx, order models.Order, reservedUntil time.Time) error {
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
/*func (r *OrdersRepository) InsertCashReservation(tx *sqlx.Tx, order models.Order, reservedUntil time.Time) error {
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
}*/

/*Updates the available_cash and reserved_cash column in accounts table */
func (r *OrdersRepository) ReserveCash(tx *sqlx.Tx, accountID uuid.UUID, amount decimal.Decimal) error {
	updateQuery := `UPDATE thyrasec.accounts SET 
                    available_cash = available_cash - $1, 
                    reserved_cash = reserved_cash + $1 
                    WHERE id = $2 AND available_cash >= $1`

	_, err := tx.Exec(updateQuery, amount, accountID)
	return err
}

func (r *OrdersRepository) ReserveAsset(tx *sqlx.Tx, accountID uuid.UUID, amount float64, assetID uuid.UUID) error {
	updateQuery := `
	UPDATE thyrasec.holdings 
	SET available_quantity = available_quantity - $1
	WHERE account_id = $2 AND asset_id = $3 AND available_quantity >= $1
`

	_, err := tx.Exec(updateQuery, float64(amount), accountID, assetID)
	return err
}

func (r *OrdersRepository) GetOrder(db *sqlx.Tx, orderID string) (*models.Order, error) {
	var order models.Order
	query := "SELECT * FROM thyrasec.orders WHERE id = $1"
	if err := db.Get(&order, query, orderID); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrdersRepository) GetOrderType(db *sqlx.Tx, id uuid.UUID) (string, error) {

	var orderTypeName string
	query := "SELECT order_type_name FROM order_types WHERE id = $1"
	if err := db.Get(&orderTypeName, query, id); err != nil {
		return "", err
	}

	return orderTypeName, nil
}

func (r *OrdersRepository) UpdateOrderStatus(tx *sqlx.Tx, orderID string, status models.OrderStatusType) error {
	query := "UPDATE thyrasec.orders SET status = $1 WHERE id = $2"
	_, err := tx.Exec(query, status, orderID)
	return err
}

/*Checks whether there is enough holdings to sell an asset */
func (r *OrdersRepository) CheckHoldings(tx *sqlx.Tx, accountID, assetID uuid.UUID) (int, error) {
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
func (r *OrdersRepository) CheckAvailableCash(tx *sqlx.Tx, accountID uuid.UUID, amount decimal.Decimal) (bool, error) {
	var availableCash decimal.Decimal
	query := `SELECT available_cash FROM thyrasec.accounts WHERE id = $1`
	err := tx.Get(&availableCash, query, accountID)
	if err != nil {
		return false, err
	}
	return availableCash.GreaterThanOrEqual(amount), nil
}

func (r *OrdersRepository) ReleaseReservation(tx *sqlx.Tx, orderID string, houseAccount string) error {
	var accountID uuid.UUID
	var amount decimal.Decimal

	// Prepare the query
	quantityQuery := `SELECT account_id, amount FROM thyrasec.cash_reservations WHERE order_id = $1`

	// Execute the query and scan the results into the variables
	err := tx.QueryRow(quantityQuery, orderID).Scan(&accountID, &amount)
	if err != nil {
		return fmt.Errorf("failed to get reservation details: %w", err)
	}

	// Update account reserved cash
	updateAccountQuery := `UPDATE thyrasec.accounts SET reserved_cash = reserved_cash - $1 WHERE id = $2`
	_, err = tx.Exec(updateAccountQuery, amount, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account reserved cash: %w", err)
	}

	updateHouseQuery := `UPDATE thyrasec.accounts SET reserved_cash = reserved_cash - $1 WHERE id = $2`
	_, err = tx.Exec(updateHouseQuery, amount, houseAccount)
	if err != nil {
		return fmt.Errorf("failed to update account reserved cash: %w", err)
	}

	// Update reservation status
	updateReservationQuery := `UPDATE thyrasec.cash_reservations SET status = 'inactive' WHERE order_id = $1`
	_, err = tx.Exec(updateReservationQuery, orderID)
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	return nil
}

func (r *OrdersRepository) UpdateAccountBalance(tx *sqlx.Tx, accountID uuid.UUID, balanceChange float64) error {
	// Define the query to update the account
	query := `UPDATE thyrasec.accounts
              SET account_balance = account_balance - $1,
                  reserved_cash = reserved_cash - $1,
                  available_cash = available_cash - $1
              WHERE id = $2`

	// Execute the query using the transaction
	if _, err := tx.Exec(query, balanceChange, accountID); err != nil {
		return err
	}

	return nil
}

func (r *OrdersRepository) InsertHolding(tx *sqlx.Tx, holding positionmodels.Holding) error {
	// Check if the holding already exists for the given account and asset
	existingHolding := positionmodels.Holding{}
	err := tx.Get(&existingHolding, "SELECT * FROM thyrasec.holdings WHERE account_id = $1 AND asset_id = $2", holding.AccountID, holding.AssetID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil {
		// Holding exists, update the quantity
		newQuantity := existingHolding.Quantity + holding.Quantity
		newAvailableQuantity := existingHolding.AvailableQuantity + holding.Quantity
		_, err := tx.Exec("UPDATE thyrasec.holdings SET quantity = $1 AND available_quantity = $2 WHERE id = $3", newQuantity, newAvailableQuantity, existingHolding.ID)
		return err
	}

	// Holding doesn't exist, insert a new row
	query := `INSERT INTO thyrasec.holdings (account_id, asset_id, quantity)
		  VALUES (:account_id, :asset_id, :quantity, :available_quantity)`
	_, err = tx.NamedExec(query, holding)
	return err
}

func (r *OrdersRepository) DeductHolding(db *sqlx.Tx, accountID uuid.UUID, assetID uuid.UUID, quantity float64) error {
	// Check if the holding exists for the given account and asset
	existingHolding := positionmodels.Holding{}
	err := db.Get(&existingHolding, "SELECT * FROM thyrasec.holdings WHERE account_id = $1 AND asset_id = $2", accountID, assetID)
	if err != nil {
		fmt.Printf("The query for fetching holding went wrong: %v", err)
		return err
	}

	// Check if the quantity to deduct is greater than the existing quantity
	if quantity > existingHolding.Quantity {
		return errors.New("insufficient holdings to deduct")
	}

	// Deduct the quantity
	newQuantity := existingHolding.Quantity - quantity
	newAvalaibleQuantity := existingHolding.AvailableQuantity - quantity
	// If the new quantity is zero, delete the holding row; otherwise, update the quantity
	if newQuantity == 0 {
		_, err = db.Exec("DELETE FROM thyrasec.holdings WHERE id = $1", existingHolding.ID)
	} else {
		_, err = db.Exec("UPDATE thyrasec.holdings SET quantity = $1, available_quantity = $2 WHERE id = $3", newQuantity, newAvalaibleQuantity, existingHolding.ID)
	}

	return err
}

func (r *OrdersRepository) UpdateOrder(db *sqlx.Tx, orderID string, settledQuantity float64, settledAmount float64, status string, tradeDate *time.Time, settlementDate *time.Time, comment string) error {
	query := `
        UPDATE thyrasec.orders
        SET settledQuantity = $1, settledAmount = $2, status = $3, trade_date = $4, settlement_date = $5, comment = $6, updated_at = NOW()
        WHERE id = $7	
        `

	_, err := db.Exec(query, settledQuantity, settledAmount, status, tradeDate, settlementDate, comment, orderID)
	return err
}

func (r *OrdersRepository) GetAssetType(db *sqlx.Tx, assetId uuid.UUID) (uuid.UUID, error) {

	var assetType uuid.UUID

	query := "SELECT asset_type_id FROM assets WHERE id = $1"
	if err := db.Get(&assetType, query, assetId); err != nil {
		return uuid.Nil, err
	}
	return assetType, nil
}

func (r *OrdersRepository) GetOrderTypeByName(tx *sqlx.Tx, name string) (uuid.UUID, error) {
	var orderType uuid.UUID
	query := "SELECT id FROM order_types WHERE order_type_name = $1"
	if err := tx.Get(&orderType, query, name); err != nil {
		return uuid.Nil, err
	}

	return orderType, nil
}

func (r *OrdersRepository) GetTransactionTypeByOrderTypeID(db *sqlx.Tx, orderTypeID uuid.UUID) (uuid.UUID, error) {
	var transactionType uuid.UUID
	query := "SELECT transaction_type_id FROM thyrasec.order_types WHERE Id = $1"
	if err := db.Get(&transactionType, query, orderTypeID); err != nil {
		return uuid.Nil, err
	}

	return transactionType, nil
}

func (r *OrdersRepository) GetAllOrders() ([]models.OrderWithDetails, error) {
	query := `
    SELECT o.*, a.account_number, asst.instrument_name, asst.instrument_type 
    FROM thyrasec.orders o
    JOIN thyrasec.accounts a ON o.account_id = a.id
    JOIN thyrasec.assets asst ON o.asset_id = asst.id
    `

	var orders []models.OrderWithDetails
	if err := r.db.Select(&orders, query); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrdersRepository) InsertOrder(tx *sqlx.Tx, order models.Order) error {
	const insertOrderQuery = `
    INSERT INTO thyrasec.orders 
    (id, account_id, asset_id, order_type, quantity, price_per_unit, total_amount, 
     status, created_at, updated_at, trade_date, settlement_date, owner_id, comment, order_number) 
    VALUES 
    (:id, :account_id, :asset_id, :order_type, :quantity, :price_per_unit, :total_amount, 
     :status, NOW(), NOW(), :trade_date, :settlement_date, :owner_id, :comment, :order_number)`

	_, err := tx.NamedExec(insertOrderQuery, order)
	return err
}

func (r *OrdersRepository) InsertCashReservation(tx *sqlx.Tx, order models.Order, reservedUntil time.Time) error {
	const insertReservationQuery = `
    INSERT INTO cash_reservations
    (order_id, reserved_until, ...)  -- replace ... with other fields as needed
    VALUES
    (:order_id, :reserved_until, ...)`

	// Prepare data for the reservation, including the order ID and reservedUntil
	reservationData := map[string]interface{}{
		"order_id":       order.ID,
		"reserved_until": reservedUntil,
		// Add other necessary data
	}

	_, err := tx.NamedExec(insertReservationQuery, reservationData)
	return err
}
