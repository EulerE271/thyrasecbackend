package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id                 uuid.UUID `json:"id" db:"id"`
	TransactionOwnerId uuid.UUID `json:"transaction_owner_id" db:"transaction_owner_id"`
	AccountOwnerId     uuid.UUID `json:"account_owner_id" db:"account_owner_id"`
	Type               uuid.UUID `json:"type" db:"type"`
	Asset1Id           uuid.UUID `json:"asset1_id" db:"asset1_id"`                 // Nullable field
	Asset2Id           uuid.UUID `json:"asset2_id" db:"asset2_id"`                 // Nullable field
	AccountAsset1Id    uuid.UUID `json:"account_asset1_id" db:"account_asset1_id"` // Nullable field
	AccountAsset2Id    uuid.UUID `json:"account_asset2_id" db:"account_asset2_id"` // Nullable field
	AmountAsset1       float64   `json:"amount_asset1" db:"amount_asset1"`         // Nullable field
	AmountAsset2       float64   `json:"amount_asset2" db:"amount_asset2"`         // Nullable field
	CreatedById        uuid.UUID `json:"created_by_id" db:"created_by_id"`
	UpdatedById        uuid.UUID `json:"updated_by_id" db:"updated_by_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	Corrected          bool      `json:"corrected" db:"corrected"`
	Canceled           bool      `json:"canceled" db:"canceled"`
	StatusTransaction  uuid.UUID `json:"status_transaction" db:"status_transaction"`
	Comment            *string   `json:"comment" db:"comment"` // Nullable field
	Asset2_currency    uuid.UUID `json:"asset2_currency" db:"asset2_currency"`
	Asset2_price       uuid.UUID `json:"asset2_price" db:"asset2_price"`
	OrderNumber        string    `json:"order_no" db:"order_no"`
	BusinessEvent      uuid.UUID `json:"business_event" db:"business_event"`
	Trade_date         time.Time `json:"trade_date" db:"trade_date"`
	Settlement_date    time.Time `json:"settlement_date" db:"settlement_date"`
}

func InitializeTransaction(createdBy, transactionType uuid.UUID, comment *string) Transaction {
	now := time.Now()
	return Transaction{
		Id:                uuid.New(), // Generate a new UUID for the transaction
		CreatedById:       createdBy,
		UpdatedById:       createdBy, // Assuming the creator is also the one performing the update
		CreatedAt:         now,
		UpdatedAt:         now,
		Corrected:         false,    // Default to false, assuming the transaction isn't corrected at the time of creation
		Canceled:          false,    // Default to false, assuming the transaction isn't canceled at the time of creation
		StatusTransaction: uuid.Nil, // Should be set to an appropriate value after initialization if needed
		Type:              transactionType,
		Comment:           comment,
		// The following fields are nullable and should be set explicitly when needed:
		// Asset1Id:
		// Asset2Id:
		// AccountAsset1Id:
		// AccountAsset2Id:
		// AmountAsset1:
		// AmountAsset2:
		// Asset2_currency:
		// Asset2_price:
		// OrderNumber:
		// BusinessEvent:
	}
}

func (t *Transaction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                   t.Id,
		"transaction_owner_id": t.TransactionOwnerId,
		"account_owner_id":     t.AccountOwnerId,
		"type":                 t.Type,
		"asset1_id":            t.Asset1Id,
		"asset2_id":            t.Asset2Id,
		"account_asset1_id":    t.AccountAsset1Id,
		"account_asset2_id":    t.AccountAsset2Id,
		"amount_asset1":        t.AmountAsset1,
		"amount_asset2":        t.AmountAsset2,
		"created_by_id":        t.CreatedById,
		"updated_by_id":        t.UpdatedById,
		"created_at":           t.CreatedAt,
		"updated_at":           t.UpdatedAt,
		"corrected":            t.Corrected,
		"canceled":             t.Canceled,
		"status_transaction":   t.StatusTransaction,
		"comment":              t.Comment,
	}
}

//Model for displaying transaction together with the owner of the transactions name
type TransactionDisplay struct {
	Transaction                        // Embed the original Transaction struct
	OwnerName                string    `json:"owner_name" db:"owner_name"`
	TypeName                 string    `json:"type_name" db:"type_name"`
	AccountAsset1AccountName string    `json:"account_asset1_account_name" db:"account_asset1_account_name"`
	AccountAsset2AccountName string    `json:"account_asset2_account_name" db:"account_asset2_account_name"`
	TransactionTypeName      string    `json:"transaction_type_name" db:"transaction_type_name"`
	StatusLabel              string    `json:"status_label" db:"status_label"`
	StatusID                 uuid.UUID `json:"status_id" db:"status_id"`
	Description              string    `json:"description" db:"description"`
	TypeID                   uuid.UUID `json:"type_id" db:"type_id"`
	AccountNumber            string    `json:"account_number" db:"account_number"`
}

//Model for returning transaction types as defined in transaction_types table
type TransactionType struct {
	Id                    int    `json:"id"`
	TypeID                string `json:"type_id" db:"type_id"`
	Transaction_type_name string `json:"transaction_type_name"`
}

//Defines the request sent to acounts service for updating balance when creating a transaction.
type BalanceRequest struct {
	AccountId uuid.UUID `json:"id" db:"id"`
	Amount    *int      `json:"amount" db:"amount"`
}
