// user_repository.go

package repositories

import (
	"context"
	"database/sql"
	"thyra/internal/users/models" // assuming this is where UserResponse is located

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers(ctx context.Context, role string) ([]models.UserResponse, error) {
	var query string
	switch role {
	case "admin":
		query = "SELECT id, username, email, customer_number FROM admins"
	case "advisor":
		query = "SELECT id, username, email, customer_number FROM partners_advisors"
	case "customer":
		query = "SELECT id, username, email, customer_number FROM customers"
	default:
		return nil, sql.ErrNoRows // Handle invalid role
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.CustomerNumber); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) GetUsernameByUUID(ctx context.Context, uuid string) (string, error) {
	var username string
	query := `SELECT full_name FROM customers WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(&username)
	return username, err
}

func (r *UserRepository) RegisterAdmin(ctx context.Context, admin models.AdminRegistrationRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO admins (username, password_hash, email, customer_number) VALUES ($1, $2, $3, $4)"
	_, err = r.db.ExecContext(ctx, query, admin.Username, string(hashedPassword), admin.Email, admin.CustomerNumber)
	return err
}

func (r *UserRepository) RegisterPartnerAdvisor(ctx context.Context, advisor models.PartnerAdvisorRegistrationRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(advisor.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO partners_advisors (username, password_hash, email, full_name, company_name, phone_number, customer_number) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = r.db.ExecContext(ctx, query, advisor.Username, string(hashedPassword), advisor.Email, advisor.FullName, advisor.CompanyName, advisor.PhoneNumber, advisor.CustomerNumber)
	return err
}

func (r *UserRepository) RegisterCustomer(ctx context.Context, customer models.CustomerRegistrationRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO customers (username, password_hash, email, full_name, address, phone_number, customer_number) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = r.db.ExecContext(ctx, query, customer.Username, string(hashedPassword), customer.Email, customer.FullName, customer.Address, customer.PhoneNumber, customer.CustomerNumber)
	return err
}
