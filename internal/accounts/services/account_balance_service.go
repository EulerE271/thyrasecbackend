package services

import (
	"context"
	repository "thyra/internal/accounts/repositories"

	"github.com/google/uuid"
)

type AccountBalanceService struct {
	repo *repository.AccountBalanceRepository
}

func NewAccountBalanceService(repo *repository.AccountBalanceRepository) *AccountBalanceService {
	return &AccountBalanceService{repo: repo}
}

// GetAggregatedValueService fetches the aggregated values for a user.
func (s *AccountBalanceService) GetAggregatedAccountValue(ctx context.Context, userId uuid.UUID) (repository.TotalValue, error) {
	// Call the repository function to get the aggregated value
	totalValue, err := s.repo.GetAggregatedValue(ctx, userId)
	if err != nil {
		return repository.TotalValue{}, err
	}

	// Process the data as needed or directly return
	return totalValue, nil
}

func (s *AccountBalanceService) GetSpecificAccountValue(ctx context.Context, accountId uuid.UUID) (repository.AccountValue, error) {
	// Call the repository function to get values for a specific account
	account, err := s.repo.GetSpecificAccountValue(ctx, accountId)
	if err != nil {
		return repository.AccountValue{}, err
	}

	// Process the data as needed or directly return
	return account, nil
}
