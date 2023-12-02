package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id                  uuid.UUID `db:"id" json:"id"`
	AccountName         string    `db:"account_name" json:"account_name"`
	AccountType         uuid.UUID `db:"account_type" json:"account_type"`
	AccountOwnerCompany bool      `db:"account_owner_company" json:"account_owner_company"`
	AccountBalance      float64   `db:"account_balance" json:"account_balance"`
	AccountCurrency     string    `db:"account_currency" json:"account_currency"`
	AccountNumber       string    `db:"account_number" json:"account_number"`
	AccountStatus       string    `db:"account_status" json:"account_status"`
	InterestRate        float64   `db:"interest_rate" json:"interest_rate"`
	OverdraftLimit      float64   `db:"overdraft_limit" json:"overdraft_limit"`
	AccountDescription  string    `db:"account_description" json:"account_description"`
	AccountHolderId     uuid.UUID `db:"account_holder_id" json:"account_holder_id"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy           uuid.UUID `db:"created_by" json:"created_by"`
	UpdatedBy           uuid.UUID `db:"updated_by" json:"updated_by"`
	AccountTypeName     string    `db:"account_type_name" json:"account_type_name"`
}

/*type TransactionMessage struct {
	TransactionID1  uuid.UUID `json:"transaction_id_1"`
	TransactionID2  uuid.UUID `json:"transaction_id_2"`
	OrderNumber     string    `json:"order_no"`
	AccountAsset2ID uuid.UUID `json:"account_asset2_id"`
	AmountAsset2    float64   `json:"amount_asset2"`
	AccountAsset1ID uuid.UUID `json:"account_asset1_id"`
	AmountAsset1    float64   `json:"amount_asset1"`
	Asset2Id        uuid.UUID `json:"asset2_id"`
	Asset1Id        uuid.UUID `json:"asset1_id"`
	ResponseTopic   string    `json:"responseTopic"` // Add this field to specify the response topic
} */

/*type ResponseMessage struct {
	TransactionID1 uuid.UUID `json:"transaction_id_1"`
	TransactionID2 uuid.UUID `json:"transaction_id_2"`
	Status         string    `json:"status"`
	Reason         string    `json:"reason,omitempty"` // Include a reason only when there is an error
} */
