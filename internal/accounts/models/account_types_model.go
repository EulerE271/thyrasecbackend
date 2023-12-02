package models

import (
	"github.com/google/uuid"
)

type AccountTypes struct {
	ID              uuid.UUID `db:"id" json:"id"`
	AccountTypeName string    `db:"account_type_name" json:"account_type_name"`
}
