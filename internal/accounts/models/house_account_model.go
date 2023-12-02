package models

import (
	"github.com/google/uuid"
)

type HouseAccount struct {
	ID uuid.UUID `db:"id" json:"id"`
}
