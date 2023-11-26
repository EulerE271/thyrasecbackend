package models

import (
	"github.com/google/uuid"
)

type Holding struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	AssetID   uuid.UUID `db:"asset_id"`
	Quantity  int       `db:"quantity"`
}
