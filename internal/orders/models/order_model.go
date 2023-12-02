package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusType string

// Define constants for OrderStatusType
const (
	StatusCreated   OrderStatusType = "created"
	StatusConfirmed OrderStatusType = "confirmed"
	StatusPending   OrderStatusType = "pending"
	StatusSettled   OrderStatusType = "settled"
	StatusExecuted  OrderStatusType = "executed"
	StatusCanceled  OrderStatusType = "canceled"
)

// Order represents an order in the database.
type Order struct {
	ID              uuid.UUID       `db:"id" json:"id"`
	AccountID       uuid.UUID       `db:"account_id" json:"account_id"`
	AssetID         uuid.UUID       `db:"asset_id" json:"asset_id"`
	Currency        uuid.UUID       `adb:"currency" json:"currency"`
	Quantity        float64         `db:"quantity" json:"quantity"`
	PricePerUnit    float64         `db:"price_per_unit" json:"price_per_unit,omitempty"` // omitempty will prevent zero value from being serialized
	TotalAmount     float64         `db:"total_amount" json:"total_amount"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	TradeDate       time.Time       `db:"trade_date" json:"trade_date"`
	SettlementDate  time.Time       `db:"settlement_date" json:"settlement_date"`
	Comment         *string         `db:"comment" json:"comment"`
	OwnerID         uuid.UUID       `db:"owner_id" json:"owner_id"`
	OrderType       uuid.UUID       `db:"order_type" json:"order_type"` // Consider making this an enum if you have a limited set of order types
	Status          OrderStatusType `db:"status" json:"status"`         // Using the defined OrderStatusType
	SettledQuantity *float64        `db:"settledquantity"`
	SettledAmount   *float64        `db:"settledamount"` //nullable
	OrderNumber     string          `db:"order_number" json:"order_number"`
}

type OrderWithDetails struct {
	Order                 // Embed the Order struct to include its fields
	AccountNumber  string `db:"account_number" json:"account_number"`
	InstrumentName string `db:"instrument_name" json:"instrument_name"`
	InstrumentType string `db:"instrument_type" json:"instrument_type"`
}
