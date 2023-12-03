// auth_repository.go

package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserCredentials(username string) (string, string, string, error) {
	userTypes := map[string]string{
		"admins":            "admin",
		"partners_advisors": "partner_advisor",
		"customers":         "customer",
	}

	var userType, storedPassword, userID string
	var err error

	for table, uType := range userTypes {
		query := `SELECT id, password_hash FROM ` + table + ` WHERE username = $1`
		err = r.db.QueryRow(query, username).Scan(&userID, &storedPassword)
		if err == nil {
			userType = uType
			break
		} else if err != sql.ErrNoRows {
			break // Exit the loop on error other than ErrNoRows
		}
	}

	return userID, storedPassword, userType, err
}
