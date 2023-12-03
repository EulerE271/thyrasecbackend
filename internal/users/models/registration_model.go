package models

// Common attributes for all user types
type BaseRegistrationRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	CustomerNumber string `json:"customer_number"`
}

// For Admin
type AdminRegistrationRequest struct {
	BaseRegistrationRequest
}

// For Partners/Advisors
type PartnerAdvisorRegistrationRequest struct {
	BaseRegistrationRequest
	FullName    string `json:"full_name"`
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
}

// For Customers
type CustomerRegistrationRequest struct {
	BaseRegistrationRequest
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}
