package services

import (
	"thyra/internal/transactions/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func InsertTransaction(database *sqlx.DB, transaction *models.Transaction, query string) (uuid.UUID, error) {
	var transactionID uuid.UUID
	tx, err := database.Beginx()
	if err != nil {
		return transactionID, err
	}

	namedStmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return transactionID, err
	}

	err = namedStmt.QueryRowx(transaction.ToMap()).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return transactionID, err
	}

	err = tx.Commit()
	if err != nil {
		return transactionID, err
	}

	return transactionID, nil
}

func InsertParentTransaction(database *sqlx.DB, transaction *models.Transaction) error {
	query := `
        INSERT INTO thyrasec.transactions (
            id, type, asset1_id, asset2_id, amount_asset1, amount_asset2, created_by_id,
            updated_by_id, created_at, updated_at, corrected, canceled, status_transaction,
            comment, transaction_owner_id, account_owner_id, account_asset1_id,
            account_asset2_id, trade_date, settlement_date, order_no
        )
        VALUES (
            :id, :type, :asset1_id, :asset2_id, :amount_asset1, :amount_asset2, :created_by_id,
            :updated_by_id, :created_at, :updated_at, :corrected, :canceled, :status_transaction,
            :comment, :transaction_owner_id, :account_owner_id, :account_asset1_id,
            :account_asset2_id, :trade_date, :settlement_date, :order_no
        )
    ` // Removed the RETURNING id part

	_, err := database.NamedExec(query, transaction)
	if err != nil {
		return err
	}

	return nil
}
