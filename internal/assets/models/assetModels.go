package models

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID             uuid.UUID `db:"id"`
	InstrumentName string    `db:"instrument_name" json:"instrument_name"`
	ISIN           string    `db:"isin"`
	Ticker         string    `db:"ticker"`
	Exchange       string    `db:"exchange"`
	Currency       string    `db:"currency"`
	InstrumentType string    `db:"instrument_type"`
	CurrentPrice   float64   `db:"current_price"`
	Volume         int64     `db:"volume"`
	Country        string    `db:"country"`
	Sector         string    `db:"sector"`
	AssetTypeId    uuid.UUID `db:"asset_type_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
