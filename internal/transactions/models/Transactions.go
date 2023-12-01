package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id                        uuid.UUID `json:"id" db:"id"`
	Type                      uuid.UUID `json:"type" db:"type"`
	AssetId                   uuid.UUID `json:"asset_id" db:"asset_id"`                 // Nullable field
	CashAmount                *float64  `json:"cash_amount" db:"cash_amount"`           // Nullable field
	AssetQuantity             *float64  `json:"asset_quantity" db:"asset_quantity"`     // Nullable field
	CashAccountId             uuid.UUID `json:"cash_account_id" db:"cash_account_id"`   // Nullable field
	AssetAccountId            uuid.UUID `json:"asset_account_id" db:"asset_account_id"` // Nullable field
	AssetType                 uuid.UUID `json:"asset_type" db:"asset_type"`
	TransactionCurrency       uuid.UUID `json:"transaction_currency" db:"transaction_currency"`
	AssetPrice                *float64  `json:"asset_price" db:"asset_price"` // Nullable field
	CreatedById               uuid.UUID `json:"created_by_id" db:"created_by_id"`
	UpdatedById               uuid.UUID `json:"updated_by_id" db:"updated_by_id"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
	Corrected                 bool      `json:"corrected" db:"corrected"`
	Canceled                  bool      `json:"canceled" db:"canceled"`
	Comment                   *string   `json:"comment" db:"comment"` // Nullable field
	TransactionOwnerId        uuid.UUID `json:"transaction_owner_id" db:"transaction_owner_id"`
	TransactionOwnerAccountId uuid.UUID `json:"transaction_owner_account_id" db:"transaction_owner_account_id"` // New field
	TradeDate                 time.Time `json:"trade_date" db:"trade_date"`
	SettlementDate            time.Time `json:"settlement_date" db:"settlement_date"`
	OrderNumber               string    `json:"order_no" db:"order_no"`
	BusinessEvent             uuid.UUID `json:"business_event" db:"business_event"`
}

func InitializeTransaction(createdBy, transactionType uuid.UUID, comment *string) Transaction {
	now := time.Now()
	return Transaction{
		Id:          uuid.New(),
		CreatedById: createdBy,
		UpdatedById: createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		Corrected:   false,
		Canceled:    false,
		Type:        transactionType,
		Comment:     comment,
		// The following fields are nullable and should be set explicitly when needed:
		// AssetId:
		// CashAmount:
		// AssetQuantity:
		// CashAccountId:
		// AssetAccountId:
		// AssetType:
		// TransactionCurrency:
		// AssetPrice:
		// TransactionOwnerId:
		// TransactionOwnerAccountId:
		// TradeDate:
		// SettlementDate:
		// OrderNumber:
		// BusinessEvent:
	}
}

func (t *Transaction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                           t.Id,
		"type":                         t.Type,
		"asset_id":                     t.AssetId,
		"cash_amount":                  t.CashAmount,
		"asset_quantity":               t.AssetQuantity,
		"cash_account_id":              t.CashAccountId,
		"asset_account_id":             t.AssetAccountId,
		"asset_type":                   t.AssetType,
		"transaction_currency":         t.TransactionCurrency,
		"asset_price":                  t.AssetPrice,
		"created_by_id":                t.CreatedById,
		"updated_by_id":                t.UpdatedById,
		"created_at":                   t.CreatedAt,
		"updated_at":                   t.UpdatedAt,
		"corrected":                    t.Corrected,
		"canceled":                     t.Canceled,
		"comment":                      t.Comment,
		"transaction_owner_id":         t.TransactionOwnerId,
		"transaction_owner_account_id": t.TransactionOwnerAccountId,
		"trade_date":                   t.TradeDate,
		"settlement_date":              t.SettlementDate,
		"order_no":                     t.OrderNumber,
		"business_event":               t.BusinessEvent,
	}
}

type TransactionDisplay struct {
	Transaction                   // Embed the updated Transaction struct
	OwnerName           string    `json:"owner_name" db:"owner_name"`
	TypeName            string    `json:"type_name" db:"type_name"`
	CashAccountName     string    `json:"cash_account_name" db:"cash_account_name"`   // Updated
	AssetAccountName    string    `json:"asset_account_name" db:"asset_account_name"` // Updated
	TransactionTypeName string    `json:"transaction_type_name" db:"transaction_type_name"`
	StatusLabel         string    `json:"status_label" db:"status_label"`
	StatusID            uuid.UUID `json:"status_id" db:"status_id"`
	Description         string    `json:"description" db:"description"`
	TypeID              uuid.UUID `json:"type_id" db:"type_id"`
	AccountNumber       string    `json:"account_number" db:"account_number"`
}

type TransactionType struct {
	Id                       int    `json:"id"`
	TypeID                   string `json:"type_id" db:"type_id"`
	Transaction_type_name    string `json:"transaction_type_name"`
	TransactionTypeShortName string `db:"trt_short_name" json:"trt_short_name"`
}

type BalanceRequest struct {
	AccountId uuid.UUID `json:"id" db:"id"`
	Amount    *int      `json:"amount" db:"amount"`
}
