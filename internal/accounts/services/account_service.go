package services

import (
	"context"
	"errors"
	"thyra/internal/accounts/models"
	repository "thyra/internal/accounts/repositories"
	"thyra/internal/accounts/utils"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, account models.Account, authUserID string) error {
	accountNumber := utils.GenerateAccountNumber(account.AccountType.String())
	return s.repo.CreateAccount(ctx, account, accountNumber, authUserID)
}

func (s *AccountService) GetAccountsByUser(ctx context.Context, userID string, authUserID string, authUserRole string) ([]models.Account, error) {
	if authUserRole != "admin" && authUserID != userID {
		return nil, errors.New("not allowed to fetch another user's accounts")
	}
	return s.repo.GetAccountsByUser(ctx, userID)
}

func (s *AccountService) GetAllAccounts(ctx context.Context, authUserRole string) ([]models.Account, error) {
	if authUserRole != "admin" {
		return nil, errors.New("only admins can fetch all accounts")
	}
	return s.repo.GetAllAccounts(ctx)
}

func (s *AccountService) GetAccountTypes(ctx context.Context, authUserRole string) ([]models.AccountTypes, error) {
	if authUserRole != "admin" {
		return nil, errors.New("only admins can fetch account types")
	}
	return s.repo.GetAccountTypes(ctx)
}

func (s *AccountService) GetHouseAccount(ctx context.Context, authUserRole string) (string, error) {
	if authUserRole != "admin" {
		return "", errors.New("only admins can fetch the house account")
	}
	return s.repo.GetHouseAccount(ctx)
}
