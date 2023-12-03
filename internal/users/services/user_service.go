// user_service.go

package services

import (
	"context"
	"thyra/internal/users/models"
	"thyra/internal/users/repositories"
	"thyra/internal/users/utils"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers(ctx context.Context, role string) ([]models.UserResponse, error) {
	// Add any business logic if needed
	return s.repo.GetAllUsers(ctx, role)
}

func (s *UserService) GetUsernameByUUID(ctx context.Context, uuid string) (string, error) {
	// Add any business logic if needed
	return s.repo.GetUsernameByUUID(ctx, uuid)
}

func (s *UserService) RegisterAdmin(ctx context.Context, admin models.AdminRegistrationRequest) error {
	admin.CustomerNumber = utils.GenerateCustomerNumber()
	return s.repo.RegisterAdmin(ctx, admin)
}

func (s *UserService) RegisterPartnerAdvisor(ctx context.Context, advisor models.PartnerAdvisorRegistrationRequest) error {
	advisor.CustomerNumber = utils.GenerateCustomerNumber()
	return s.repo.RegisterPartnerAdvisor(ctx, advisor)
}

func (s *UserService) RegisterCustomer(ctx context.Context, customer models.CustomerRegistrationRequest) error {
	customer.CustomerNumber = utils.GenerateCustomerNumber()
	return s.repo.RegisterCustomer(ctx, customer)
}
