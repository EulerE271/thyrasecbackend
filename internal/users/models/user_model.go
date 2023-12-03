package models

type UserResponse struct {
	UUID           string `json:"uuid"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	CustomerNumber string `json:"customer_number"`
	// Add any other fields you want to be returned in the API response
}
