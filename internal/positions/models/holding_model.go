package models

import (
	"github.com/google/uuid"
)

type Holding struct {
	ID                uuid.UUID `db:"id"`
	AccountID         uuid.UUID `db:"account_id"`
	AssetID           uuid.UUID `db:"asset_id"`
	Quantity          float64   `db:"quantity"`
	AvailableQuantity float64   `db:"available_quantity"`
}
