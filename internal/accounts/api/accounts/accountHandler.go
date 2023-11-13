// handlers/account_handler.go
package handlers

import (
	"net/http"
	"thyra/internal/accounts/models" // Import your Account model
	accountno "thyra/internal/accounts/services/accountnumber"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateAccountHandler(c *gin.Context) {
	sqlxDB := c.MustGet("db").(*sqlx.DB)

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountNumber := accountno.GenerateAccountNumber(account.AccountType.String())
	authUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserID, exists = c.Get("userID")
	if !exists {
		return
	}

	// Insert the new account into the database
	_, err := sqlxDB.NamedExec(`
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func GetAccountsByUser(c *gin.Context) {
	// Extract the target userId from the route parameter.
	targetUserID := c.Param("userId")

	// Extract the authenticated user's ID and role from context.
	authUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType") // Use "userType" instead of "userRole"
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Ensure the authenticated user is allowed to access the target user's accounts.
	if authUserRole != "admin" && authUserID.(string) != targetUserID { // Cast authUserID to string for comparison
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to fetch another user's accounts"})
		return
	}

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	// Fetch all accounts for the targetUserID (account_holder_id)
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
    account_types
ON
    accounts.account_type = account_types.id
WHERE
    accounts.account_holder_id = $1
`

	var accounts []models.Account
	if err := sqlxDB.Select(&accounts, query, targetUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No accounts were found"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func GetAllAccounts(c *gin.Context) {
	// Extract the authenticated user's ID and role from context.
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Ensure the authenticated user is an admin.
	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to fetch all accounts"})
		return
	}

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	// Since it's an admin, fetch all accounts from the database
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
	account_types
  ON
	accounts.account_type = account_types.id
  `
	var accounts []models.Account
	if err := sqlxDB.Select(&accounts, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch accounts", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func GetAccountTypes(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Ensure the authenticated user is an admin.
	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to fetch all accounts"})
		return
	}

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	query := "SELECT * FROM account_types"
	var accountTypes []models.AccountTypes
	if err := sqlxDB.Select(&accountTypes, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message:": "Failed to fetch account types", "error": err.Error()})
	}

	c.JSON(http.StatusOK, accountTypes)
}

func GetHouseAccount(c *gin.Context) {

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUserRole, exists := c.Get("userType")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Ensure the authenticated user is an admin.
	if authUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must be an admin to fetch all accounts"})
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}
	sqlxDB, ok := db.(*sqlx.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database type assertion failed"})
		return
	}

	// Define the query to join accounts with account_types to find the house account
	query := `
        SELECT a.id
        FROM accounts a
        INNER JOIN account_types at ON a.account_type = at.id
        WHERE at.account_type_name = $1;
    `

	var accountID string
	// Execute the query with "House" as the account type name we are looking for
	err := sqlxDB.Get(&accountID, query, "House")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching house account", "details": err.Error()})
		return
	}

	// Respond with the found house account ID
	c.JSON(http.StatusOK, gin.H{"house_account_id": accountID})
}
