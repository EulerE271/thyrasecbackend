package repositories

import (
	"thyra/internal/transactions/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	InsertTransaction(transaction *models.Transaction) error
	UpdateAccountBalance(accountID uuid.UUID, newBalance float64, availableBalance float64) error
	GetAccountBalance(accountID uuid.UUID) (float64, error)
	GetAccountAvailableBalance(accountID uuid.UUID) (float64, error)
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) InsertTransaction(transaction *models.Transaction) error {
	query := `
        INSERT INTO thyrasec.transactions (
            id, type, asset_id, cash_amount, asset_quantity, cash_account_id, 
            asset_account_id, asset_type, transaction_currency, asset_price, 
            created_by_id, updated_by_id, created_at, updated_at, corrected, canceled,
            comment, transaction_owner_id, transaction_owner_account_id, trade_date, 
            settlement_date, order_no, business_event
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
            $17, $18, $19, $20, $21, $22, $23
        )`
	_, err := r.db.Exec(query,
		transaction.Id, transaction.Type, transaction.AssetId, transaction.CashAmount,
		transaction.AssetQuantity, transaction.CashAccountId, transaction.AssetAccountId,
		transaction.AssetType, transaction.TransactionCurrency, transaction.AssetPrice,
		transaction.CreatedById, transaction.UpdatedById, transaction.CreatedAt,
		transaction.UpdatedAt, transaction.Corrected, transaction.Canceled,
		transaction.Comment, transaction.TransactionOwnerId, transaction.TransactionOwnerAccountId,
		transaction.TradeDate, transaction.SettlementDate, transaction.OrderNumber,
		transaction.BusinessEvent)
	return err
}

func (r *transactionRepository) UpdateAccountBalance(accountID uuid.UUID, newBalance, availableBalance float64) error {
	_, err := r.db.Exec("UPDATE accounts SET account_balance = $1, available_cash = $2, updated_at = $3 WHERE id = $4", newBalance, availableBalance, time.Now(), accountID)
	return err
}

func (r *transactionRepository) GetAccountBalance(accountID uuid.UUID) (float64, error) {
	var balance float64
	err := r.db.QueryRow("SELECT account_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *transactionRepository) GetAccountAvailableBalance(accountID uuid.UUID) (float64, error) {
	var availableBalance float64
	err := r.db.QueryRow("SELECT available_cash FROM accounts WHERE id = $1", accountID).Scan(&availableBalance)
	if err != nil {
		return 0, err
	}

	return availableBalance, nil
}

func (r *transactionRepository) GetAccountAvailableIinstrument(accountID uuid.UUID, instrumentID uuid.UUID) (float64, error) {
	var availableQuantity float64
	err := r.db.QueryRow("SELECT quantity FROM holdings WHERE account_id = $1 AND asset_id = $2", accountID, instrumentID).Scan(&availableQuantity)
	if err != nil {
		return 0, err
	}

	return availableQuantity, nil
}
