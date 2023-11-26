package models

import (
	"time"

	"github.com/google/uuid"
)

type Instrument struct {
	ID             uuid.UUID `db:"id" json:"id"`
	InstrumentName string    `db:"instrument_name" json:"instrument_name"`
	ISIN           string    `db:"isin" json:"isin"`
	Ticker         string    `db:"ticker" json:"ticker"`
	Exchange       string    `db:"exchange" json:"exchange"`
	Currency       string    `db:"currency" json:"currency"`
	InstrumentType string    `db:"instrument_type" json:"instrument_type"`
	CurrentPrice   float64   `db:"current_price" json:"current_price"`
	Volume         int64     `db:"volume" json:"volume"`
	Country        string    `db:"country" json:"country"`
	Sector         string    `db:"sector" json:"sector"`
	AssetTypeID    uuid.UUID `db:"asset_type_id" json:"asset_type_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

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
	ID           uuid.UUID       `db:"id" json:"id"`
	AccountID    uuid.UUID       `db:"account_id" json:"account_id"`
	AssetID      uuid.UUID       `db:"asset_id" json:"asset_id"`
	OrderType    string          `db:"order_type" json:"order_type"` // Consider making this an enum if you have a limited set of order types
	Quantity     int             `db:"quantity" json:"quantity"`
	PricePerUnit float64         `db:"price_per_unit" json:"price_per_unit,omitempty"` // omitempty will prevent zero value from being serialized
	TotalAmount  float64         `db:"total_amount" json:"total_amount"`
	Status       OrderStatusType `db:"status" json:"status"` // Using the defined OrderStatusType
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

type OrderWithDetails struct {
	Order                 // Embed the Order struct to include its fields
	AccountNumber  string `db:"account_number" json:"account_number"`
	InstrumentName string `db:"instrument_name" json:"instrument_name"`
	InstrumentType string `db:"instrument_type" json:"instrument_type"`
}
