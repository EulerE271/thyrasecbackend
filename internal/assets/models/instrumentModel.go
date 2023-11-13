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

type Assets struct {
	Id          uuid.UUID `db:"id" json:"id"`
	TypeName    string    `db:"type_name" json:"type_name"`
	Description string    `db:"description" json:"description"`
}
