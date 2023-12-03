package repository

import (
	"context" // Added import for logging
	"thyra/internal/accounts/models"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account models.Account, accountNumber string, authUserID string) error {
	_, err := r.db.NamedExecContext(ctx, `
        INSERT INTO thyrasec.accounts (account_name, account_type, account_owner_company, account_balance, account_currency,
            account_number, account_status, interest_rate, overdraft_limit,
            account_description, account_holder_id, created_by, updated_by)
        VALUES (:account_name, :account_type, :account_owner_company, :account_balance, :account_currency,
            :account_number, :account_status, :interest_rate, :overdraft_limit,
            :account_description, :account_holder_id, :created_by, :updated_by)
    `, map[string]interface{}{
		"account_name":          account.AccountName,
		"account_type":          account.AccountType,
		"account_owner_company": account.AccountOwnerCompany,
		"account_balance":       account.AccountBalance,
		"account_currency":      account.AccountCurrency,
		"account_number":        accountNumber,

		"account_status":      account.AccountStatus,
		"interest_rate":       account.InterestRate,
		"overdraft_limit":     account.OverdraftLimit,
		"account_description": account.AccountDescription,
		"account_holder_id":   account.AccountHolderId,
		"created_by":          authUserID,
		"updated_by":          authUserID,
	})

	return err
}

func (r *AccountRepository) GetAccountsByUser(ctx context.Context, userID string) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT
	accounts.id,
	accounts.account_name,
	accounts.account_balance,
	accounts.account_currency,
	accounts.account_number,
	accounts.account_status,
	accounts.interest_rate,
	accounts.overdraft_limit,
	accounts.account_description,
	accounts.account_holder_id,
	accounts.created_at,
	accounts.updated_at,
	accounts.created_by,
	accounts.updated_by,
	account_types.account_type_name
FROM
	accounts
INNER JOIN
	account_types ON accounts.account_type = account_types.id
WHERE
	accounts.account_holder_id = $1
`

	err := r.db.SelectContext(ctx, &accounts, query, userID)
	return accounts, err
}

func (r *AccountRepository) GetAllAccounts(ctx context.Context) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT
        accounts.id,
        accounts.account_name,
        accounts.account_balance,
        accounts.account_currency,
        accounts.account_number,
        accounts.account_status,
        accounts.interest_rate,
        accounts.overdraft_limit,
        accounts.account_description,
        accounts.account_holder_id,
        accounts.created_at,
        accounts.updated_at,
        accounts.created_by,
        accounts.updated_by,
        account_types.account_type_name
    FROM
        accounts
    INNER JOIN
        account_types ON accounts.account_type = account_types.id
    `
	err := r.db.SelectContext(ctx, &accounts, query)
	return accounts, err
}

func (r *AccountRepository) GetAccountTypes(ctx context.Context) ([]models.AccountTypes, error) {
	var accountTypes []models.AccountTypes
	query := "SELECT * FROM account_types"
	err := r.db.SelectContext(ctx, &accountTypes, query)
	return accountTypes, err
}

func (r *AccountRepository) GetHouseAccount(ctx context.Context) (string, error) {
	var accountID string
	query := `
        SELECT a.id
        FROM accounts a
        INNER JOIN account_types at ON a.account_type = at.id
        WHERE at.account_type_name = $1;
    `
	err := r.db.GetContext(ctx, &accountID, query, "House")
	return accountID, err
}
