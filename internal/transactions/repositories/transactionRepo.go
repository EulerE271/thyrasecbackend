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
            id, type, asset1_id, asset2_id, amount_asset1, amount_asset2, created_by_id,
            updated_by_id, created_at, updated_at, corrected, canceled, status_transaction,
            comment, transaction_owner_id, account_owner_id, account_asset1_id,
            account_asset2_id, trade_date, settlement_date, order_no
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
        )`
	_, err := r.db.Exec(query,
		transaction.Id, transaction.Type, transaction.Asset1Id, transaction.Asset2Id,
		transaction.AmountAsset1, transaction.AmountAsset2, transaction.CreatedById,
		transaction.UpdatedById, transaction.CreatedAt, transaction.UpdatedAt,
		transaction.Corrected, transaction.Canceled, transaction.StatusTransaction,
		transaction.Comment, transaction.TransactionOwnerId, transaction.AccountOwnerId,
		transaction.AccountAsset1Id, transaction.AccountAsset2Id,
		transaction.Trade_date, transaction.Settlement_date, transaction.OrderNumber)
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